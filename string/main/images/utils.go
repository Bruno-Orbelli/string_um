package images

import (
	"fmt"
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

	jpegBytes, err := gocv.IMEncode(".jpg", img)
	if err != nil {
		return nil, err
	}

	// Detect and recognize faces
	faces, err := rec.Recognize(jpegBytes.GetBytes())
	if err != nil {
		return nil, err
	} else if len(faces) == 0 {
		return nil, fmt.Errorf("no faces in image or unrecognized face")
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
func NormalizeEncoding(encoding []float32) []float64 {
	var norm float64
	for _, v := range encoding {
		norm += float64(v * v)
	}
	norm = math.Sqrt(norm)

	normalized := make([]float64, len(encoding))
	for i, v := range encoding {
		normalized[i] = float64(v) / norm
	}
	return normalized
}
