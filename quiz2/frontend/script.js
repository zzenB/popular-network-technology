const fileInput = document.getElementById("file-input");
const uploadButton = document.getElementById("upload-button");
const progressSection = document.getElementById("progress-section");
const progressFill = document.getElementById("progress-fill");
const statusMessage = document.getElementById("status-message");
const downloadSection = document.getElementById("download-section");
const downloadLink = document.getElementById("download-link");
const MAX_FILE_SIZE = 75 * 1024 * 1024; // 75 MB

// Backend url
const backendUrl = "https://pdf.rifaldoagustinus.com/upload"; // Replace with actual backend URL

const fileInfo = document.createElement("p");
fileInfo.style.marginTop = "10px";
fileInfo.style.color = "#666";
document.getElementById("upload-section").appendChild(fileInfo);

fileInput.addEventListener("change", () => {
  const file = fileInput.files[0];
  if (file) {
    const fileSizeMB = (file.size / (1024 * 1024)).toFixed(2);
    fileInfo.textContent = `Selected file: ${file.name} (${fileSizeMB} MB)`;
  } else {
    fileInfo.textContent = "";
  }
});


// Event listener for the upload button
uploadButton.addEventListener("click", () => {
  try {
    uploadButton.disabled = true;
    uploadButton.textContent = "Uploading...";

    const file = fileInput.files[0];

    if (!file) {
      alert("Please select a PDF file to upload!");
      return;
    }

    if (file.type !== "application/pdf") {
      alert("Please upload a valid PDF file.");
      return;
    }

    if (file.size > MAX_FILE_SIZE) {
      alert(`File size exceeds the 75 MB limit. Please select a smaller file.`);
      return;
    }

    uploadFile(file);
  } finally {
    uploadButton.disabled = false;
    uploadButton.textContent = "Upload and Compress";
  }
});

// Upload file function
async function uploadFile(file) {
  progressSection.style.display = "block";
  progressFill.style.width = "0%";
  statusMessage.textContent = "Uploading...";

  const formData = new FormData();
  formData.append("file", file);

  try {
    const response = await fetch(backendUrl, {
      method: "POST",
      body: formData,
    });

    // Log the response for debugging
    console.log('Response status:', response.status);
    const responseText = await response.text();
    console.log('Response text:', responseText);

    if (!response.ok) {
      const errorDetails = await response.json();
      throw new Error(errorDetails.error || "File upload failed!");
    }

    const result = await response.json();
    const fileId = result.fileId;

    progressFill.style.width = "100%";
    downloadSection.style.display = "block";

    const originalFileName = file.name;
    const compressedFileName = `compressed_${originalFileName}`;

    downloadLink.href = `${backendUrl.replace('/upload', '/download')}/${fileId}`;
    downloadLink.textContent = `${compressedFileName}`;
  } catch (error) {
    progressFill.style.width = "0%";
    statusMessage.textContent = "An error occurred during upload.";

    if (error.message.includes('NetworkError')) {
      alert("Network error. Please check your connection and try again.");
    }
    else {
      alert("Error: " + error.message);
    }

    // Log error for debugging
    console.error('Upload error:', error);
  } finally {
    uploadButton.disabled = false;
    uploadButton.textContent = "Upload and Compress";
  }
}
