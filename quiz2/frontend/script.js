// script.js

// Select DOM elements
const fileInput = document.getElementById("file-input");
const uploadButton = document.getElementById("upload-button");
const progressSection = document.getElementById("progress-section");
const progressFill = document.getElementById("progress-fill");
const statusMessage = document.getElementById("status-message");
const downloadSection = document.getElementById("download-section");
const downloadLink = document.getElementById("download-link");
const MAX_FILE_SIZE = 10 * 1024 * 1024; // 10 MB

// Backend endpoint
const backendUrl = "https://pdf.rifaldoagustinus.com/upload"; // Replace with actual backend URL

// Event listener for the upload button
uploadButton.addEventListener("click", () => {
  const file = fileInput.files[0];

  if (!file) {
    alert("Please select a PDF file to upload!");
    return;
  }

  // Check file type
  if (file.type !== "application/pdf") {
    alert("Please upload a valid PDF file.");
    return;
  }

  // Check file size
  if (file.size > MAX_FILE_SIZE) {
    alert(`File size exceeds the 10 MB limit. Please select a smaller file.`);
    return;
  }

  // Proceed with the upload
  uploadFile(file);
});

// Upload file function
async function uploadFile(file) {
  // Show progress section and reset progress bar
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

    if (!response.ok) {
      const errorDetails = await response.json();
      throw new Error(errorDetails.error || "File upload failed!");
    }

    const result = await response.json();
    const fileId = result.fileId;

    progressFill.style.width = "100%";
    statusMessage.textContent = "Compression complete!";
    downloadSection.style.display = "block";
    downloadLink.href = `${backendUrl.replace('/upload', '/download')}/${fileId}`;
    downloadLink.textContent = "Download Compressed PDF";
  } catch (error) {
    progressFill.style.width = "0%";
    statusMessage.textContent = "An error occurred during upload.";
    alert("Error: " + error.message);
  }
}
