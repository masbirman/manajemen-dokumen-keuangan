<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from "vue";

const emit = defineEmits<{
  "pdf-created": [file: File];
  close: [];
}>();

const videoRef = ref<HTMLVideoElement | null>(null);
const canvasRef = ref<HTMLCanvasElement | null>(null);
const stream = ref<MediaStream | null>(null);
const capturedImages = ref<string[]>([]);
const isCapturing = ref(false);
const error = ref("");

const isMobile = computed(() => {
  return /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(
    navigator.userAgent
  );
});

const startCamera = async () => {
  try {
    error.value = "";
    stream.value = await navigator.mediaDevices.getUserMedia({
      video: {
        facingMode: "environment",
        width: { ideal: 1920 },
        height: { ideal: 1080 },
      },
    });
    if (videoRef.value) {
      videoRef.value.srcObject = stream.value;
    }
  } catch (e) {
    error.value =
      "Tidak dapat mengakses kamera. Pastikan izin kamera diberikan.";
    console.error(e);
  }
};

const stopCamera = () => {
  if (stream.value) {
    stream.value.getTracks().forEach((track) => track.stop());
    stream.value = null;
  }
};

const captureImage = () => {
  if (!videoRef.value || !canvasRef.value) return;

  const video = videoRef.value;
  const canvas = canvasRef.value;
  canvas.width = video.videoWidth;
  canvas.height = video.videoHeight;

  const ctx = canvas.getContext("2d");
  if (ctx) {
    ctx.drawImage(video, 0, 0);
    const imageData = canvas.toDataURL("image/jpeg", 0.8);
    capturedImages.value.push(imageData);
  }
};

const removeImage = (index: number) => {
  capturedImages.value.splice(index, 1);
};

const createPDF = async () => {
  if (capturedImages.value.length === 0) return;

  isCapturing.value = true;
  try {
    // Dynamic import jsPDF
    const { jsPDF } = await import("jspdf");
    const pdf = new jsPDF({
      orientation: "portrait",
      unit: "mm",
      format: "a4",
    });

    for (let i = 0; i < capturedImages.value.length; i++) {
      if (i > 0) pdf.addPage();

      const img = new Image();
      img.src = capturedImages.value[i];
      await new Promise((resolve) => (img.onload = resolve));

      const pageWidth = pdf.internal.pageSize.getWidth();
      const pageHeight = pdf.internal.pageSize.getHeight();
      const imgRatio = img.width / img.height;
      const pageRatio = pageWidth / pageHeight;

      let imgWidth, imgHeight;
      if (imgRatio > pageRatio) {
        imgWidth = pageWidth - 20;
        imgHeight = imgWidth / imgRatio;
      } else {
        imgHeight = pageHeight - 20;
        imgWidth = imgHeight * imgRatio;
      }

      const x = (pageWidth - imgWidth) / 2;
      const y = (pageHeight - imgHeight) / 2;
      pdf.addImage(capturedImages.value[i], "JPEG", x, y, imgWidth, imgHeight);
    }

    const pdfBlob = pdf.output("blob");
    const file = new File([pdfBlob], `scan_${Date.now()}.pdf`, {
      type: "application/pdf",
    });
    emit("pdf-created", file);
    stopCamera();
  } catch (e) {
    error.value = "Gagal membuat PDF";
    console.error(e);
  } finally {
    isCapturing.value = false;
  }
};

onMounted(() => {
  if (isMobile.value) startCamera();
});

onUnmounted(() => {
  stopCamera();
});
</script>

<template>
  <div v-if="isMobile" class="fixed inset-0 bg-black z-50 flex flex-col">
    <div class="bg-gray-900 text-white p-4 flex justify-between items-center">
      <h2 class="font-semibold">Scan Dokumen</h2>
      <button @click="$emit('close')" class="text-2xl">&times;</button>
    </div>

    <div v-if="error" class="bg-red-500 text-white p-3 text-center">
      {{ error }}
    </div>

    <div class="flex-1 relative">
      <video
        ref="videoRef"
        autoplay
        playsinline
        class="w-full h-full object-cover"
      ></video>
      <canvas ref="canvasRef" class="hidden"></canvas>

      <div class="absolute bottom-4 left-0 right-0 flex justify-center gap-4">
        <button
          @click="captureImage"
          class="w-16 h-16 bg-white rounded-full border-4 border-gray-300 flex items-center justify-center"
        >
          <div class="w-12 h-12 bg-red-500 rounded-full"></div>
        </button>
      </div>
    </div>

    <div v-if="capturedImages.length > 0" class="bg-gray-900 p-4">
      <div class="flex gap-2 overflow-x-auto pb-2">
        <div
          v-for="(img, index) in capturedImages"
          :key="index"
          class="relative flex-shrink-0"
        >
          <img :src="img" class="w-16 h-20 object-cover rounded" />
          <button
            @click="removeImage(index)"
            class="absolute -top-2 -right-2 w-5 h-5 bg-red-500 text-white rounded-full text-xs"
          >
            &times;
          </button>
        </div>
      </div>
      <div class="flex gap-2 mt-3">
        <span class="text-white text-sm"
          >{{ capturedImages.length }} halaman</span
        >
        <button
          @click="createPDF"
          :disabled="isCapturing"
          class="flex-1 bg-blue-600 text-white py-2 rounded-lg disabled:opacity-50"
        >
          {{ isCapturing ? "Membuat PDF..." : "Buat PDF" }}
        </button>
      </div>
    </div>
  </div>

  <div v-else class="text-center p-8 bg-gray-100 rounded-lg">
    <p class="text-gray-500">
      📱 Fitur scan dokumen hanya tersedia di perangkat mobile
    </p>
  </div>
</template>
