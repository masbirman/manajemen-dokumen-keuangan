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
const containerRef = ref<HTMLDivElement | null>(null);
const stream = ref<MediaStream | null>(null);
const capturedImages = ref<string[]>([]);
const isCapturing = ref(false);
const error = ref("");
const isVideoReady = ref(false);
const scanMode = ref<"single" | "multi">("single");
const compressionQuality = ref(0.7);
const estimatedSize = ref(0);

const MAX_PDF_SIZE = 300 * 1024; // 300KB
const CROP_PADDING = 16; // padding from edge in pixels (matching inset-4 = 1rem = 16px)

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

// Enhancement modes
type EnhanceMode = "original" | "enhanced" | "bw" | "sharp";
const enhanceMode = ref<EnhanceMode>("enhanced");

// Apply image enhancement based on mode
const applyEnhancement = (
  ctx: CanvasRenderingContext2D,
  mode: EnhanceMode,
  canvas: HTMLCanvasElement
) => {
  const imageData = ctx.getImageData(0, 0, canvas.width, canvas.height);
  const data = imageData.data;

  switch (mode) {
    case "enhanced":
      // High contrast + brightness boost for documents
      for (let i = 0; i < data.length; i += 4) {
        // Increase contrast
        data[i] = Math.min(255, Math.max(0, (data[i] - 128) * 1.3 + 128 + 15));     // R
        data[i + 1] = Math.min(255, Math.max(0, (data[i + 1] - 128) * 1.3 + 128 + 15)); // G
        data[i + 2] = Math.min(255, Math.max(0, (data[i + 2] - 128) * 1.3 + 128 + 15)); // B
      }
      break;

    case "bw":
      // Convert to black & white with threshold
      for (let i = 0; i < data.length; i += 4) {
        const gray = 0.299 * data[i] + 0.587 * data[i + 1] + 0.114 * data[i + 2];
        // Apply threshold for cleaner B&W
        const bw = gray > 140 ? 255 : 0;
        data[i] = bw;
        data[i + 1] = bw;
        data[i + 2] = bw;
      }
      break;

    case "sharp":
      // Grayscale + sharpen for text documents
      const grayData = new Uint8ClampedArray(data.length);
      for (let i = 0; i < data.length; i += 4) {
        const gray = 0.299 * data[i] + 0.587 * data[i + 1] + 0.114 * data[i + 2];
        // Increase contrast on grayscale
        const enhanced = Math.min(255, Math.max(0, (gray - 128) * 1.4 + 128 + 20));
        grayData[i] = enhanced;
        grayData[i + 1] = enhanced;
        grayData[i + 2] = enhanced;
        grayData[i + 3] = data[i + 3];
      }
      for (let i = 0; i < data.length; i++) {
        data[i] = grayData[i];
      }
      break;

    case "original":
    default:
      // No enhancement
      break;
  }

  ctx.putImageData(imageData, 0, 0);
};

// Compress and enhance image
const compressImage = async (
  dataUrl: string,
  maxWidth: number = 1200,
  quality: number = 0.7,
  mode: EnhanceMode = "enhanced"
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
        ctx.drawImage(img, 0, 0, width, height);
        applyEnhancement(ctx, mode, canvas);
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
    const base64Data = img.split(",")[1] || "";
    totalSize += (base64Data.length * 3) / 4;
  });
  totalSize += images.length * 5000;
  return totalSize;
};

// Capture and crop image to the guide frame area
const captureImage = async () => {
  if (!videoRef.value || !canvasRef.value || !containerRef.value || !isVideoReady.value) {
    error.value = "Kamera belum siap, tunggu sebentar...";
    return;
  }

  const video = videoRef.value;
  const canvas = canvasRef.value;
  const container = containerRef.value;
  
  const videoWidth = video.videoWidth;
  const videoHeight = video.videoHeight;
  
  // Get container dimensions (visible area)
  const containerRect = container.getBoundingClientRect();
  const displayWidth = containerRect.width;
  const displayHeight = containerRect.height;
  
  // Calculate the crop area based on the guide frame (inset-4 = 16px padding)
  const paddingPercent = CROP_PADDING / Math.min(displayWidth, displayHeight);
  
  // Since video is object-cover, calculate the actual visible portion
  const videoAspect = videoWidth / videoHeight;
  const containerAspect = displayWidth / displayHeight;
  
  let sourceX = 0;
  let sourceY = 0;
  let sourceWidth = videoWidth;
  let sourceHeight = videoHeight;
  
  if (videoAspect > containerAspect) {
    // Video is wider - crop sides
    sourceWidth = videoHeight * containerAspect;
    sourceX = (videoWidth - sourceWidth) / 2;
  } else {
    // Video is taller - crop top/bottom
    sourceHeight = videoWidth / containerAspect;
    sourceY = (videoHeight - sourceHeight) / 2;
  }
  
  // Apply the guide frame padding (crop to the inner area)
  const cropPadding = Math.min(sourceWidth, sourceHeight) * paddingPercent;
  sourceX += cropPadding;
  sourceY += cropPadding;
  sourceWidth -= cropPadding * 2;
  sourceHeight -= cropPadding * 2;
  
  // Set canvas to cropped dimensions (A4-like aspect ratio for documents)
  const targetAspect = 210 / 297; // A4 aspect ratio
  let cropWidth = sourceWidth;
  let cropHeight = sourceHeight;
  
  // Adjust to closer to A4 aspect if possible
  const currentAspect = cropWidth / cropHeight;
  if (currentAspect > targetAspect) {
    cropWidth = cropHeight * targetAspect;
    sourceX += (sourceWidth - cropWidth) / 2;
  } else {
    cropHeight = cropWidth / targetAspect;
    sourceY += (sourceHeight - cropHeight) / 2;
  }
  
  // Set output canvas size
  canvas.width = Math.round(cropWidth);
  canvas.height = Math.round(cropHeight);

  const ctx = canvas.getContext("2d");
  if (ctx) {
    // Draw cropped portion
    ctx.drawImage(
      video,
      sourceX, sourceY, cropWidth, cropHeight,  // Source (crop area)
      0, 0, canvas.width, canvas.height         // Destination (full canvas)
    );
    
    const rawImage = canvas.toDataURL("image/jpeg", 0.9);
    
    // Compress and enhance image
    const compressedImage = await compressImage(rawImage, 1200, compressionQuality.value, enhanceMode.value);
    
    if (scanMode.value === "single") {
      capturedImages.value = [compressedImage];
    } else {
      capturedImages.value.push(compressedImage);
    }
    
    // Update estimated size
    estimatedSize.value = estimatePdfSize(capturedImages.value);
    
    // Auto-adjust quality if too large
    if (estimatedSize.value > MAX_PDF_SIZE && compressionQuality.value > 0.3) {
      compressionQuality.value -= 0.1;
      const recompressed = await Promise.all(
        capturedImages.value.map((img) => compressImage(img, 1000, compressionQuality.value, enhanceMode.value))
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

      const compressedImages = await Promise.all(
        capturedImages.value.map((img) => compressImage(img, 1000, quality, enhanceMode.value))
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
        <p class="text-xs text-gray-400">Max 300KB • Auto-crop aktif</p>
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
        📄 Single
      </button>
      <button
        @click="scanMode = 'multi'; capturedImages = []"
        class="flex-1 py-2 rounded-lg text-sm font-medium transition-colors"
        :class="scanMode === 'multi' 
          ? 'bg-blue-600 text-white' 
          : 'bg-gray-700 text-gray-300'"
      >
        📚 Multi
      </button>
    </div>

    <!-- Enhancement Mode Toggle -->
    <div class="bg-gray-800 px-4 pb-2 flex gap-1">
      <button
        @click="enhanceMode = 'original'"
        class="flex-1 py-1.5 rounded text-xs font-medium transition-colors"
        :class="enhanceMode === 'original' 
          ? 'bg-green-600 text-white' 
          : 'bg-gray-700 text-gray-400'"
      >
        Original
      </button>
      <button
        @click="enhanceMode = 'enhanced'"
        class="flex-1 py-1.5 rounded text-xs font-medium transition-colors"
        :class="enhanceMode === 'enhanced' 
          ? 'bg-green-600 text-white' 
          : 'bg-gray-700 text-gray-400'"
      >
        ✨ Enhanced
      </button>
      <button
        @click="enhanceMode = 'bw'"
        class="flex-1 py-1.5 rounded text-xs font-medium transition-colors"
        :class="enhanceMode === 'bw' 
          ? 'bg-green-600 text-white' 
          : 'bg-gray-700 text-gray-400'"
      >
        ⬛ B&W
      </button>
      <button
        @click="enhanceMode = 'sharp'"
        class="flex-1 py-1.5 rounded text-xs font-medium transition-colors"
        :class="enhanceMode === 'sharp' 
          ? 'bg-green-600 text-white' 
          : 'bg-gray-700 text-gray-400'"
      >
        🔍 Sharp
      </button>
    </div>

    <!-- Error -->
    <div v-if="error" class="bg-red-500 text-white p-3 text-center text-sm">
      {{ error }}
    </div>

    <!-- Camera View -->
    <div ref="containerRef" class="flex-1 relative overflow-hidden">
      <video
        ref="videoRef"
        autoplay
        playsinline
        class="w-full h-full object-cover"
      ></video>

      <canvas ref="canvasRef" class="hidden"></canvas>

      <!-- Document Guide Frame - this defines crop area -->
      <div class="absolute inset-4 border-2 border-white/70 rounded-lg pointer-events-none bg-transparent">
        <!-- Corner markers -->
        <div class="absolute -top-0.5 -left-0.5 w-8 h-8 border-t-4 border-l-4 border-green-400 rounded-tl-lg"></div>
        <div class="absolute -top-0.5 -right-0.5 w-8 h-8 border-t-4 border-r-4 border-green-400 rounded-tr-lg"></div>
        <div class="absolute -bottom-0.5 -left-0.5 w-8 h-8 border-b-4 border-l-4 border-green-400 rounded-bl-lg"></div>
        <div class="absolute -bottom-0.5 -right-0.5 w-8 h-8 border-b-4 border-r-4 border-green-400 rounded-br-lg"></div>
        
        <!-- Instruction text -->
        <div class="absolute inset-0 flex items-center justify-center pointer-events-none">
          <p class="bg-black/50 text-white px-3 py-1.5 rounded-lg text-sm backdrop-blur-sm">
            Posisikan dokumen di dalam frame
          </p>
        </div>
      </div>
      
      <!-- Dark overlay outside crop area -->
      <div class="absolute inset-0 pointer-events-none">
        <div class="absolute top-0 left-0 right-0 h-4 bg-black/40"></div>
        <div class="absolute bottom-0 left-0 right-0 h-4 bg-black/40"></div>
        <div class="absolute top-4 bottom-4 left-0 w-4 bg-black/40"></div>
        <div class="absolute top-4 bottom-4 right-0 w-4 bg-black/40"></div>
      </div>

      <!-- Capture Button -->
      <div class="absolute bottom-6 left-0 right-0 flex justify-center">
        <button
          @click="captureImage"
          :disabled="!isVideoReady"
          class="w-18 h-18 bg-white rounded-full border-4 border-green-400 flex items-center justify-center disabled:opacity-50 transition-all active:scale-95 shadow-lg"
          style="width: 72px; height: 72px;"
        >
          <div class="w-14 h-14 bg-green-500 rounded-full" style="width: 56px; height: 56px;"></div>
        </button>
      </div>

      <!-- Single Mode: Preview Overlay -->
      <div 
        v-if="scanMode === 'single' && capturedImages.length > 0"
        class="absolute inset-0 bg-black/90 flex items-center justify-center p-4"
      >
        <div class="bg-white rounded-xl p-4 max-w-sm w-full shadow-2xl">
          <p class="text-sm text-gray-600 text-center mb-3 font-medium">Preview Hasil Crop</p>
          <img :src="capturedImages[0]" class="w-full rounded-lg mb-4 border" />
          <p class="text-sm text-gray-600 text-center mb-4">
            Ukuran: ~{{ formatSize(estimatedSize) }}
          </p>
          <div class="flex gap-2">
            <button
              @click="capturedImages = []"
              class="flex-1 py-2.5 bg-gray-200 text-gray-700 rounded-lg font-medium"
            >
              ↻ Ulangi
            </button>
            <button
              @click="createPDF"
              :disabled="isCapturing"
              class="flex-1 py-2.5 bg-blue-600 text-white rounded-lg font-medium disabled:opacity-50"
            >
              {{ isCapturing ? "Proses..." : "✓ Simpan" }}
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
          <img :src="img" class="w-14 h-18 object-cover rounded border-2 border-gray-700" style="height: 72px;" />
          <span class="absolute bottom-0 left-0 right-0 bg-black/70 text-white text-xs text-center py-0.5">
            {{ index + 1 }}
          </span>
          <button
            @click="removeImage(index)"
            class="absolute -top-1 -right-1 w-5 h-5 bg-red-500 text-white rounded-full text-xs flex items-center justify-center shadow"
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
          {{ isCapturing ? "Membuat PDF..." : "✓ Buat PDF" }}
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
