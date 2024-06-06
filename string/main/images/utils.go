package images

import (
	"fmt"
	"image"
	"math"

	"github.com/Kagami/go-face"
	"gocv.io/x/gocv"
)

// Function to capture an image from the camera
func CaptureImage() (gocv.Mat, error) {
	// Open the default camera
	webcam, err := gocv.VideoCaptureDevice(0)
	if err != nil {
		return gocv.Mat{}, fmt.Errorf("error opening video capture device: %v", err)
	}
	defer webcam.Close()

	// Create a Mat to store the image
	img := gocv.NewMat()
	defer img.Close()

	// Read a frame from the camera
	if ok := webcam.Read(&img); !ok {
		return gocv.Mat{}, fmt.Errorf("cannot read from webcam")
	}

	return img.Clone(), nil
}

// DetectAlignAndRescale detects faces in an image, aligns them and rescales them
func DetectAlignAndRescale(img gocv.Mat) (gocv.Mat, error) {
	// Load the face detector and shape predictor
	classifier := gocv.NewCascadeClassifier()
	if !classifier.Load("models/haarcascade_frontalface_default.xml") {
		return gocv.Mat{}, fmt.Errorf("error loading face detector model")
	}
	defer classifier.Close()

	// Detect faces
	rects := classifier.DetectMultiScale(img)
	if len(rects) == 0 {
		return gocv.Mat{}, fmt.Errorf("no faces detected")
	}

	// Align face (simple implementation)
	alignedFace := img.Region(rects[0])

	// Rescale the face
	resizedFace := gocv.NewMat()
	gocv.Resize(alignedFace, &resizedFace, image.Point{X: 160, Y: 160}, 0, 0, gocv.InterpolationLinear)
	defer resizedFace.Close()

	gocv.IMWrite("image.jpg", resizedFace)

	return resizedFace.Clone(), nil
}

func GrayScaleAndEqualizeHist(img gocv.Mat) gocv.Mat {
	grayScale := gocv.NewMatWithSize(img.Rows(), img.Cols(), gocv.MatTypeCV8UC1)
	defer grayScale.Close()
	grayScale.ConvertTo(&grayScale, gocv.MatTypeCV8UC1)
	img.ConvertTo(&img, gocv.MatTypeCV8UC1)

	// Convert the image to grayscale
	gocv.CvtColor(img, &grayScale, gocv.ColorRGBToGray)

	equalized := gocv.NewMatWithSize(img.Rows(), img.Cols(), gocv.MatTypeCV8UC1)
	defer equalized.Close()

	// Equalize the histogram
	gocv.EqualizeHist(grayScale, &equalized)

	// Write the image to a file
	gocv.IMWrite("image.jpg", equalized)

	return equalized.Clone()
}

// Function to extract features from an image
func ExtractFeatures(img gocv.Mat) ([]float32, error) {
	// Initialize the face recognizer
	rec, err := face.NewRecognizer("models")
	if err != nil {
		return nil, err
	}
	defer rec.Close()

	// Read the image
	if img.Empty() {
		return nil, fmt.Errorf("could not read image")
	}
	defer img.Close()

	bgrImg := gocv.NewMatWithSize(img.Rows(), img.Cols(), gocv.MatTypeCV8UC3)
	gocv.CvtColor(img, &bgrImg, gocv.ColorGrayToBGR)
	jpegEncode, err := gocv.IMEncode(".jpg", bgrImg)
	gocv.IMWrite("image.jpg", bgrImg)
	if err != nil {
		return nil, err
	}

	// Detect and recognize faces
	faces, err := rec.Recognize(jpegEncode.GetBytes())
	if err != nil {
		return nil, err
	} else if len(faces) == 0 {
		return nil, fmt.Errorf("no faces in image")
	}

	// Return the embedding of the first detected face
	return faces[0].Descriptor[:], nil
}

// Function to average multiple face encodings
func AverageEncodings(encodings [][]float32) []float32 {
	numEncodings := len(encodings)
	if numEncodings == 0 {
		return nil
	}

	avgEncoding := make([]float32, len(encodings[0]))
	for _, encoding := range encodings {
		for i, value := range encoding {
			avgEncoding[i] += value
		}
	}

	for i := range avgEncoding {
		avgEncoding[i] /= float32(numEncodings)
	}

	return avgEncoding
}

// Function to normalize an encoding
func NormalizeEncoding(encoding []float32) []float32 {
	norm := float32(0.0)
	for _, v := range encoding {
		norm += float32(v * v)
	}
	norm = float32(math.Sqrt(float64(norm)))

	normalized := make([]float32, len(encoding))
	for i, v := range encoding {
		normalized[i] = float32(v) / norm
	}

	return normalized
}

func RoundEncoding(encodedImage []float32) []float32 {
	roundedImage := make([]float32, len(encodedImage))
	for i, value := range encodedImage {
		roundedImage[i] = float32(math.Round(float64(value)*1000)) / 1000 // Round to three decimal places
	}

	return roundedImage
}
