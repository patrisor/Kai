# Kai: Your Personal AI-Powered Operating System Assistant

Kai is your intelligent virtual assistant designed to seamlessly integrate with your computer, transforming the way you manage your digital tasks. Kai is highly organized, efficient, reliable, and trustworthy. Whether it's managing your system processes, executing commands, or providing valuable insights, Kai is here to make your life easier and more convenient. The longer you use Kai, the better it gets at understanding and assisting you, offering a truly personalized experience.

## Features

- **Operating System Management**: Kai can manage the OS of any computer, handling memory, processes, software, and hardware.
- **Command Execution**: Executes system-specific shell commands to achieve user requests.
- **Iterative Problem Solving**: Analyzes command outputs and iteratively comes up with solutions until the task is successfully completed.
- **Learning Capability**: Learns about you over time, enhancing the experience the longer you use it.

## Getting Started

### Prerequisites

Ensure you have the following installed before starting:

- **Go 1.16 or later**: [Install Go](https://golang.org/dl/)
- **Git**: [Install Git](https://git-scm.com/)
- **pkg-config**: Required for compiling with certain dependencies
  - **macOS**: `brew install pkg-config`
  - **Linux**: `sudo apt-get install pkg-config`
  - **Windows**: Use MSYS2 or another package manager
- **PortAudio**: Required for audio processing
  - **macOS**: `brew install portaudio`
  - **Linux**: `sudo apt-get install libportaudio2 libportaudio-dev`
  - **Windows**: Manually install or use MSYS2

### API Setup (Must be completed before running the application)

1. **Obtain a Gemini API Key**:
   - Follow the instructions in the [Obtaining a Gemini API Key](#obtaining-a-gemini-api-key) section to get your API key.

2. **Set Up Google Cloud API Credentials**:
   - Follow the instructions in the [Setting Up Google Cloud API Credentials](#setting-up-google-cloud-api-credentials) section to download and place the required service account JSON file in the correct directory.


Ensure you have the following installed before starting:

- **Go 1.16 or later**: [Install Go](https://golang.org/dl/)
- **Git**: [Install Git](https://git-scm.com/)
- **pkg-config**: Required for compiling with certain dependencies
  - **macOS**: `brew install pkg-config`
  - **Linux**: `sudo apt-get install pkg-config`
  - **Windows**: Use MSYS2 or another package manager
- **PortAudio**: Required for audio processing
  - **macOS**: `brew install portaudio`
  - **Linux**: `sudo apt-get install libportaudio2 libportaudio-dev`
  - **Windows**: Manually install or use MSYS2

### Installation

1. **Clone the Repository**:
   ```sh
   git clone https://github.com/patrisor/kai.git
   cd kai
   ```

2. **Complete API Setup**:
   - Make sure you have completed the API setup, including obtaining your Gemini API key and setting up your Google Cloud API credentials.

### Build and Run

1. **Build and Run**:
   ```sh
   go build && ./kai
   ```

**Important**: Do not build or run the application until you have completed the API setup, including obtaining your Gemini API key and setting up your Google Cloud API credentials.

1. **Install or Update Go**:
   - **macOS**: `brew install go` (or `brew upgrade go` to update)
   - **Windows/Linux**: [Download and install Go](https://golang.org/dl/)

2. **Clone the Repository**:
   ```sh
   git clone https://github.com/patrisor/kai.git
   cd kai
   ```

3. **Build and Run**:
   ```sh
   go build && ./kai
   ```

### Obtaining a Gemini API Key

Before you can start using Kai, you'll need to obtain a Gemini API key. Follow these steps:

1. **Sign Up for the Gemini API**:
   - Visit the [Gemini API Portal](https://api.gemini.com/).
   - Create an account if you don’t already have one.

2. **Generate an API Key**:
   - Once logged in, navigate to the API section.
   - Create a new API key, ensuring you have the necessary permissions for the features you plan to use.

3. **Store Your API Key**:
   - Keep your API key secure. You'll need to input it when you first run Kai.
   - Optionally, you can set the API key as an environment variable named `GEMINI_API_KEY`:
     ```sh
     export GEMINI_API_KEY=your_api_key_here
     ```

When you run Kai for the first time, you will be prompted to enter your Gemini API key. Make sure to have it ready.

### Setting Up Google Cloud API Credentials

To enable speech-to-text and text-to-speech features:

1. **Download the Service Account File**:
   - Go to the [Google Cloud Console](https://console.cloud.google.com/).
   - Navigate to **IAM & Admin** > **Service Accounts**.
   - Select or create a service account with the necessary permissions (e.g., "Cloud Speech-to-Text API User").
   - Click **Manage Keys** > **Add Key** > **Create New Key**.
   - Choose **JSON** and click **Create**. The file will download automatically.

2. **Move the File to the `.config` Directory**:
   - Place the downloaded JSON file in the `.config` directory within your `kai` project folder:
     ```sh
     mv ~/Downloads/gen-lang-client-<identifier>.json ./kai/.config/
     ```

3. **Run the Application**:
   ```sh
   go build && ./kai
   ```

Kai will automatically detect the service account file and set up the necessary environment variable.

## Contributing

Contributions are welcome! If you have suggestions or find any issues, feel free to open an issue or submit a pull request.

## License

This project is licensed under a Proprietary License. See the [LICENSE](./LICENSE) file for details.
