<script setup>
try {
  const xhr = new XMLHttpRequest();
  xhr.open("POST", "http://localhost:9091/upload", true);

  xhr.upload.onprogress = (event) => {
    if (event.lengthComputable) {
      const totalBytes = event.total || 0;
      const uploadedBytes = event.loaded || 0;

      selectedFiles.value.forEach((file) => {
        const uploadProgress = Math.round((uploadedBytes / totalBytes) * 100);
        file.uploadProgress = uploadProgress;
        console.log(`Dosya: ${file.name} - Yüzde: ${uploadProgress}%`);
      });
    }
  };

  xhr.onload = async () => {
    const response = JSON.parse(xhr.responseText);
    if (response.otp) {
      otpMessage.value = response.otp;
      currentStep.value = 3; // OTP gösterim adımına geç
    } else {
      otpMessage.value = "Hata: OTP alınamadı.";
    }
  };

  xhr.onerror = () => {
    console.error("Yükleme hatası oluştu");
    otpMessage.value = "Sunucu ile iletişimde bir sorun var.";
  };

  xhr.send(formData);

} catch (error) {
  console.error("Hata:", error);
  otpMessage.value = "Sunucu ile iletişimde bir sorun var.";
}

</script>