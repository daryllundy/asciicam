🖼️ asciicam

[![GitHub](https://img.shields.io/badge/GitHub-Primary-181717?logo=github)](https://github.com/daryllundy/asciicam) [![GitLab](https://img.shields.io/badge/GitLab-Mirror-FCA121?logo=gitlab)](https://gitlab.com/daryllundy/asciicam)


Display your webcam as ASCII art directly in the terminal.
Compatible with macOS and Linux.

⸻

📦 Requirements
	•	Go 1.17 or newer
	•	OpenCV (via gocv) for webcam support

⸻

🛠️ Installation (macOS Example)
	1.	Install OpenCV:

brew install opencv

	2.	Ensure pkg-config can find OpenCV:

pkg-config --cflags --libs opencv4

	3.	Build the app:

go build

⚠️ If you encounter build issues, verify your OpenCV installation and pkg-config setup.

⸻

▶️ Usage

Basic:

asciicam

Options:

Flag	Description
-dev=1	Use a specific camera device (default: 0)
-width=80 -height=60	Set output dimensions (auto-adjusts to terminal size by default)
-camWidth=640 -camHeight=480	Set camera input resolution (default: 1280x720)
-zoom=2	Set zoom level: 1 = 25%, 2 = 50%, 3 = 75%, 4 = 100% (default)
-color="#00ff00"	Output in monochrome with a custom color
-ansi=true	Use ANSI color output (good for modern terminals)
-gen=true -sample bgdata/	Generate background sample data (for greenscreen mode)
-greenscreen=true	Enable virtual greenscreen (requires sample data)
-threshold=0.12	Set greenscreen sensitivity threshold (works with -greenscreen)
-fps=true	Display frames-per-second counter



⸻

🖼️ Screenshots

ANSI Mode:

ASCII Mode:
