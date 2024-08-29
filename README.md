# Kai: Your Personal AI-Powered Operating System Assistant

Kai is your intelligent virtual assistant designed to seamlessly integrate with your computer, transforming the way you manage your digital tasks. Kai is highly organized, efficient, reliable, and trustworthy. Whether it's managing your system processes, executing commands, or providing valuable insights, Kai is here to make your life easier and more convenient. The longer you use Kai, the better it gets at understanding and assisting you, offering a truly personalized experience.

## Features

- **Operating System Management**: Kai can manage the OS of any computer, handling memory, processes, software, and hardware.
- **Command Execution**: Executes system-specific shell commands to achieve user requests.
- **Iterative Problem Solving**: Analyzes command outputs and iteratively comes up with solutions until the task is successfully completed.
- **Learning Capability**: Learns about you over time, enhancing the experience the longer you use it.

## Getting Started

To begin using Kai, follow the steps below to set up your environment.

### Prerequisites

Ensure you have the following installed before starting:

- **Go 1.16 or later**:
  - **macOS**: `brew install go` (or `brew upgrade go` to update)
  - **Windows/Linux**: [Download and install Go](https://golang.org/dl/)
- **Git**: [Install Git](https://git-scm.com/)
- **pkg-config**: Required for compiling with certain dependencies
  - **macOS**: `brew install pkg-config`
  - **Linux**: `sudo apt-get install pkg-config`
  - **Windows**: Use MSYS2 or another package manager
- **PortAudio**: Required for audio processing
  - **macOS**: `brew install portaudio`
  - **Linux**: `sudo apt-get install libportaudio2 libportaudio-dev`
  - **Windows**: Manually install or use MSYS2

### Step 1: API Setup

Kai requires two APIs to function properly: the Gemini API and Google Cloud API. Ensure you complete these setups before running the application.

#### Obtaining a Gemini API Key

1. **Sign Up for the Gemini API**:
   - Visit the [Gemini API Portal](https://api.gemini.com/) and create an account if you don’t already have one.

2. **Generate an API Key**:
   - Navigate to the API section and create a new API key with the necessary permissions.

3. **Store Your API Key**:
   - Keep your API key secure. You'll need to input it when you first run Kai.

#### Setting Up Google Cloud API Credentials

To enable speech-to-text and text-to-speech features:

1. **Download the Service Account File**:
   - Go to the [Google Cloud Console](https://console.cloud.google.com/).
   - Navigate to **IAM & Admin** > **Service Accounts**.
   - Select or create a service account with the necessary permissions for both the **Cloud Text-to-Speech API** and **Cloud Speech-to-Text API**.
   - Click **Manage Keys** > **Add Key** > **Create New Key**.
   - Choose **JSON** and click **Create** to download the file.

2. **Move the File to the `.config` Directory**:
   - Place the downloaded JSON file in the `.config` directory within your `kai` project folder.

### Step 2: Installation

After completing the API setup, proceed with the following steps to install and run Kai.

1. **Clone the Repository**:
   ```sh
   git clone https://github.com/patrisor/kai.git
   cd kai
   ```

2. **Build and Run**:
   ```sh
   go build && ./kai
   ```

**Important**: Do not attempt to build or run the application until you have completed the API setup, including obtaining your Gemini API key and setting up your Google Cloud API credentials.

## Contributing

Contributions are welcome! If you have suggestions or find any issues, feel free to open an issue or submit a pull request.

## License

This project is licensed under a Proprietary License. See the [LICENSE](./LICENSE) file for details.
