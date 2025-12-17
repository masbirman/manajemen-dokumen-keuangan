<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from "vue";

const props = defineProps<{
  nilai?: number;
  tanggal?: string;
}>();

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
const isVideoReady = ref(false);
const scanMode = ref<"single" | "multi">("single");
const compressionQuality = ref(0.7);
const estimatedSize = ref(0);

const MAX_PDF_SIZE = 300 * 1024; // 300KB

const isMobile = computed(() => {
  return /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(
    navigator.userAgent
  );
});

// Generate PDF filename based on nilai and tanggal
const generateFileName = () => {
  const nilai = props.nilai || 0;
  const tanggal = props.tanggal || new Date().toISOString().split("T")[0];
  return `DOK_${nilai}_${tanggal}.pdf`;
};

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
      videoRef.value.onloadedmetadata = () => {
        isVideoReady.value = true;
      };
    }
  } catch (e) {
    error.value = "Tidak dapat mengakses kamera. Pastikan izin kamera diberikan.";
    console.error(e);
  }
};

const stopCamera = () => {
  if (stream.value) {
    stream.value.getTracks().forEach((track) => track.stop());
    stream.value = null;
  }
};

// Compress image to target size
const compressImage = async (
  dataUrl: string,
  maxWidth: number = 1200,
  quality: number = 0.7
): Promise<string> => {
  return new Promise((resolve) => {
    const img = new Image();
    img.onload = () => {
      const canvas = document.createElement("canvas");
      let width = img.width;
      let height = img.height;

      // Scale down if too large
      if (width > maxWidth) {
        height = (height * maxWidth) / width;
        width = maxWidth;
      }

      canvas.width = width;
      canvas.height = height;

      const ctx = canvas.getContext("2d");
      if (ctx) {
        // Apply slight contrast enhancement
        ctx.filter = "contrast(1.1) brightness(1.05)";
        ctx.drawImage(img, 0, 0, width, height);
      }

      resolve(canvas.toDataURL("image/jpeg", quality));
    };
    img.src = dataUrl;
  });
};

// Estimate PDF size from images
const estimatePdfSize = (images: string[]): number => {
  let totalSize = 0;
  images.forEach((img) => {
    // Base64 is ~4/3 of binary size, plus PDF overhead
    const base64Data = img.split(",")[1] || "";
    totalSize += (base64Data.length * 3) / 4;
  });
  // Add PDF structure overhead (~5KB per page)
  totalSize += images.length * 5000;
  return totalSize;
};

const captureImage = async () => {
  if (!videoRef.value || !canvasRef.value || !isVideoReady.value) {
    error.value = "Kamera belum siap, tunggu sebentar...";
    return;
  }

  const video = videoRef.value;
  const canvas = canvasRef.value;
  canvas.width = video.videoWidth;
  canvas.height = video.videoHeight;

  const ctx = canvas.getContext("2d");
  if (ctx) {
    ctx.drawImage(video, 0, 0);
    const rawImage = canvas.toDataURL("image/jpeg", 0.9);
    
    // Compress image
    const compressedImage = await compressImage(rawImage, 1200, compressionQuality.value);
    
    if (scanMode.value === "single") {
      // Single mode: replace existing image
      capturedImages.value = [compressedImage];
    } else {
      // Multi mode: add to list
      capturedImages.value.push(compressedImage);
    }
    
    // Update estimated size
    estimatedSize.value = estimatePdfSize(capturedImages.value);
    
    // Auto-adjust quality if too large
    if (estimatedSize.value > MAX_PDF_SIZE && compressionQuality.value > 0.3) {
      compressionQuality.value -= 0.1;
      // Recompress all images with lower quality
      const recompressed = await Promise.all(
        capturedImages.value.map((img) => compressImage(img, 1000, compressionQuality.value))
      );
      capturedImages.value = recompressed;
      estimatedSize.value = estimatePdfSize(capturedImages.value);
    }
  }
};

const removeImage = (index: number) => {
  capturedImages.value.splice(index, 1);
  estimatedSize.value = estimatePdfSize(capturedImages.value);
};

const createPDF = async () => {
  if (capturedImages.value.length === 0) return;

  isCapturing.value = true;
  error.value = "";
  
  try {
    const { jsPDF } = await import("jspdf");
    
    // Start with reasonable quality
    let quality = compressionQuality.value;
    let pdfBlob: Blob | null = null;
    let attempts = 0;
    const maxAttempts = 5;
    
    while (attempts < maxAttempts) {
      const pdf = new jsPDF({
        orientation: "portrait",
        unit: "mm",
        format: "a4",
      });

      // Compress images with current quality
      const compressedImages = await Promise.all(
        capturedImages.value.map((img) => compressImage(img, 1000, quality))
      );

      for (let i = 0; i < compressedImages.length; i++) {
        if (i > 0) pdf.addPage();

        const img = new Image();
        img.src = compressedImages[i];
        await new Promise((resolve) => (img.onload = resolve));

        const pageWidth = pdf.internal.pageSize.getWidth();
        const pageHeight = pdf.internal.pageSize.getHeight();
        const imgRatio = img.width / img.height;
        const pageRatio = pageWidth / pageHeight;

        let imgWidth, imgHeight;
        if (imgRatio > pageRatio) {
          imgWidth = pageWidth - 10;
          imgHeight = imgWidth / imgRatio;
        } else {
          imgHeight = pageHeight - 10;
          imgWidth = imgHeight * imgRatio;
        }

        const x = (pageWidth - imgWidth) / 2;
        const y = (pageHeight - imgHeight) / 2;
        pdf.addImage(compressedImages[i], "JPEG", x, y, imgWidth, imgHeight);
      }

      pdfBlob = pdf.output("blob");
      
      if (pdfBlob.size <= MAX_PDF_SIZE) {
        break;
      }
      
      // Reduce quality and try again
      quality -= 0.1;
      if (quality < 0.2) quality = 0.2;
      attempts++;
    }
    
    if (!pdfBlob) {
      throw new Error("Failed to create PDF");
    }
    
    if (pdfBlob.size > MAX_PDF_SIZE) {
      error.value = `PDF terlalu besar (${(pdfBlob.size / 1024).toFixed(0)}KB). Coba kurangi jumlah halaman.`;
      return;
    }

    const fileName = generateFileName();
    const file = new File([pdfBlob], fileName, {
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

const formatSize = (bytes: number): string => {
  return (bytes / 1024).toFixed(0) + "KB";
};

onMounted(() => {
  if (isMobile.value) {
    startCamera();
  }
});

onUnmounted(() => {
  stopCamera();
});
</script>

<template>
  <div v-if="isMobile" class="fixed inset-0 bg-black z-50 flex flex-col">
    <!-- Header -->
    <div class="bg-gray-900 text-white p-4 flex justify-between items-center">
      <div>
        <h2 class="font-semibold">Scan Dokumen</h2>
        <p class="text-xs text-gray-400">Max 300KB</p>
      </div>
      <button @click="$emit('close')" class="text-2xl">&times;</button>
    </div>

    <!-- Mode Toggle -->
    <div class="bg-gray-800 px-4 py-2 flex gap-2">
      <button
        @click="scanMode = 'single'; capturedImages = []"
        class="flex-1 py-2 rounded-lg text-sm font-medium transition-colors"
        :class="scanMode === 'single' 
          ? 'bg-blue-600 text-white' 
          : 'bg-gray-700 text-gray-300'"
      >
        📄 Single Page
      </button>
      <button
        @click="scanMode = 'multi'; capturedImages = []"
        class="flex-1 py-2 rounded-lg text-sm font-medium transition-colors"
        :class="scanMode === 'multi' 
          ? 'bg-blue-600 text-white' 
          : 'bg-gray-700 text-gray-300'"
      >
        📚 Multi Page
      </button>
    </div>

    <!-- Error -->
    <div v-if="error" class="bg-red-500 text-white p-3 text-center text-sm">
      {{ error }}
    </div>

    <!-- Camera View -->
    <div class="flex-1 relative overflow-hidden">
      <video
        ref="videoRef"
        autoplay
        playsinline
        class="w-full h-full object-cover"
      ></video>

      <canvas ref="canvasRef" class="hidden"></canvas>

      <!-- Document Guide Frame -->
      <div class="absolute inset-4 border-2 border-white/60 rounded-lg pointer-events-none">
        <div class="absolute top-0 left-0 w-6 h-6 border-t-4 border-l-4 border-green-400 rounded-tl-lg"></div>
        <div class="absolute top-0 right-0 w-6 h-6 border-t-4 border-r-4 border-green-400 rounded-tr-lg"></div>
        <div class="absolute bottom-0 left-0 w-6 h-6 border-b-4 border-l-4 border-green-400 rounded-bl-lg"></div>
        <div class="absolute bottom-0 right-0 w-6 h-6 border-b-4 border-r-4 border-green-400 rounded-br-lg"></div>
      </div>

      <!-- Capture Button -->
      <div class="absolute bottom-4 left-0 right-0 flex justify-center">
        <button
          @click="captureImage"
          :disabled="!isVideoReady"
          class="w-16 h-16 bg-white rounded-full border-4 border-green-400 flex items-center justify-center disabled:opacity-50 transition-all active:scale-95 shadow-lg"
        >
          <div class="w-12 h-12 bg-green-500 rounded-full"></div>
        </button>
      </div>

      <!-- Single Mode: Preview Overlay -->
      <div 
        v-if="scanMode === 'single' && capturedImages.length > 0"
        class="absolute inset-0 bg-black/80 flex items-center justify-center p-4"
      >
        <div class="bg-white rounded-xl p-4 max-w-sm w-full">
          <img :src="capturedImages[0]" class="w-full rounded-lg mb-4" />
          <p class="text-sm text-gray-600 text-center mb-4">
            Ukuran: ~{{ formatSize(estimatedSize) }}
          </p>
          <div class="flex gap-2">
            <button
              @click="capturedImages = []"
              class="flex-1 py-2 bg-gray-200 text-gray-700 rounded-lg font-medium"
            >
              Ulangi
            </button>
            <button
              @click="createPDF"
              :disabled="isCapturing"
              class="flex-1 py-2 bg-blue-600 text-white rounded-lg font-medium disabled:opacity-50"
            >
              {{ isCapturing ? "Proses..." : "Simpan" }}
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Multi Mode: Thumbnail Strip & Actions -->
    <div v-if="scanMode === 'multi' && capturedImages.length > 0" class="bg-gray-900 p-4">
      <div class="flex gap-2 overflow-x-auto pb-2">
        <div
          v-for="(img, index) in capturedImages"
          :key="index"
          class="relative flex-shrink-0"
        >
          <img :src="img" class="w-14 h-18 object-cover rounded border-2 border-gray-700" />
          <span class="absolute bottom-0 left-0 right-0 bg-black/70 text-white text-xs text-center">
            {{ index + 1 }}
          </span>
          <button
            @click="removeImage(index)"
            class="absolute -top-1 -right-1 w-5 h-5 bg-red-500 text-white rounded-full text-xs flex items-center justify-center"
          >
            ×
          </button>
        </div>
      </div>
      
      <div class="flex items-center gap-3 mt-3">
        <div class="text-white text-sm">
          <span class="font-medium">{{ capturedImages.length }}</span> hal
          <span class="text-gray-400 ml-2">~{{ formatSize(estimatedSize) }}</span>
          <span 
            v-if="estimatedSize > MAX_PDF_SIZE" 
            class="text-red-400 ml-1"
          >
            (melebihi batas)
          </span>
        </div>
        <button
          @click="createPDF"
          :disabled="isCapturing || capturedImages.length === 0"
          class="flex-1 bg-blue-600 text-white py-2.5 rounded-lg font-medium disabled:opacity-50"
        >
          {{ isCapturing ? "Membuat PDF..." : "Buat PDF" }}
        </button>
      </div>
    </div>
  </div>

  <!-- Desktop Message -->
  <div v-else class="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4">
    <div class="bg-white rounded-xl shadow-xl max-w-sm w-full p-6 text-center">
      <div class="w-16 h-16 bg-blue-100 text-blue-600 rounded-full flex items-center justify-center mx-auto mb-4 text-3xl">
        📱
      </div>
      <h3 class="text-lg font-bold text-gray-800 mb-2">Gunakan Smartphone</h3>
      <p class="text-gray-600 mb-6">
        Fitur scan kamera hanya tersedia jika diakses melalui browser HP / Tablet.
      </p>
      <button
        @click="$emit('close')"
        class="w-full py-2.5 bg-gray-100 text-gray-700 font-medium rounded-lg hover:bg-gray-200 transition-colors"
      >
        Mengerti
      </button>
    </div>
  </div>
</template>
