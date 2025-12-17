<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch } from "vue";

const props = defineProps<{
  nilai?: number;
  tanggal?: string;
}>();

const emit = defineEmits<{
  "pdf-created": [file: File];
  close: [];
}>();

// OpenCV types
declare const cv: {
  Mat: new () => Mat;
  MatVector: new () => MatVector;
  Size: new (width: number, height: number) => Size;
  cvtColor: (src: Mat, dst: Mat, code: number) => void;
  GaussianBlur: (src: Mat, dst: Mat, ksize: Size, sigmaX: number) => void;
  Canny: (src: Mat, dst: Mat, threshold1: number, threshold2: number) => void;
  findContours: (src: Mat, contours: MatVector, hierarchy: Mat, mode: number, method: number) => void;
  contourArea: (contour: Mat) => number;
  arcLength: (curve: Mat, closed: boolean) => number;
  approxPolyDP: (curve: Mat, approxCurve: Mat, epsilon: number, closed: boolean) => void;
  getPerspectiveTransform: (src: Mat, dst: Mat) => Mat;
  warpPerspective: (src: Mat, dst: Mat, M: Mat, dsize: Size) => void;
  matFromImageData: (imageData: ImageData) => Mat;
  COLOR_RGBA2GRAY: number;
  RETR_EXTERNAL: number;
  CHAIN_APPROX_SIMPLE: number;
  matFromArray: (rows: number, cols: number, type: number, array: number[]) => Mat;
  CV_32FC2: number;
};

interface Mat {
  delete: () => void;
  rows: number;
  cols: number;
  data: Uint8Array;
  data32F: Float32Array;
  type: () => number;
  size: () => { width: number; height: number };
}

interface MatVector {
  delete: () => void;
  size: () => number;
  get: (index: number) => Mat;
}

interface Size {
  width: number;
  height: number;
}

const videoRef = ref<HTMLVideoElement | null>(null);
const canvasRef = ref<HTMLCanvasElement | null>(null);
const overlayCanvasRef = ref<HTMLCanvasElement | null>(null);
const stream = ref<MediaStream | null>(null);
const capturedImages = ref<string[]>([]);
const isCapturing = ref(false);
const error = ref("");
const isOpenCVLoaded = ref(false);
const isVideoReady = ref(false);
const detectedCorners = ref<{ x: number; y: number }[]>([]);
const isDetecting = ref(false);

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

// Load OpenCV.js from CDN
const loadOpenCV = (): Promise<void> => {
  return new Promise((resolve, reject) => {
    if (typeof cv !== "undefined" && cv.Mat) {
      resolve();
      return;
    }

    const script = document.createElement("script");
    script.src = "https://docs.opencv.org/4.9.0/opencv.js";
    script.async = true;

    script.onload = () => {
      const checkReady = () => {
        if (typeof cv !== "undefined" && cv.Mat) {
          resolve();
        } else {
          setTimeout(checkReady, 100);
        }
      };
      checkReady();
    };

    script.onerror = () => reject(new Error("Failed to load OpenCV.js"));
    document.head.appendChild(script);
  });
};

// Initialize scanner
const initScanner = async () => {
  try {
    await loadOpenCV();
    isOpenCVLoaded.value = true;
  } catch (e) {
    console.error("Failed to load OpenCV:", e);
    // Continue without edge detection
    isOpenCVLoaded.value = false;
  }
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
        if (isOpenCVLoaded.value) {
          startEdgeDetection();
        }
      };
    }
  } catch (e) {
    error.value =
      "Tidak dapat mengakses kamera. Pastikan izin kamera diberikan.";
    console.error(e);
  }
};

const stopCamera = () => {
  isDetecting.value = false;
  if (stream.value) {
    stream.value.getTracks().forEach((track) => track.stop());
    stream.value = null;
  }
};

// Edge detection using OpenCV
let animationFrameId: number | null = null;

const orderPoints = (pts: { x: number; y: number }[]): { x: number; y: number }[] => {
  // Order: top-left, top-right, bottom-right, bottom-left
  const sorted = [...pts].sort((a, b) => a.y - b.y);
  const top = sorted.slice(0, 2).sort((a, b) => a.x - b.x);
  const bottom = sorted.slice(2, 4).sort((a, b) => a.x - b.x);
  return [top[0], top[1], bottom[1], bottom[0]];
};

const startEdgeDetection = () => {
  if (!videoRef.value || !overlayCanvasRef.value || !isOpenCVLoaded.value) return;

  isDetecting.value = true;
  const video = videoRef.value;
  const overlayCanvas = overlayCanvasRef.value;
  const ctx = overlayCanvas.getContext("2d");

  const detectFrame = () => {
    if (!isDetecting.value || !ctx || !isOpenCVLoaded.value) return;

    try {
      // Get video dimensions
      const width = video.videoWidth;
      const height = video.videoHeight;
      
      if (width === 0 || height === 0) {
        animationFrameId = requestAnimationFrame(detectFrame);
        return;
      }

      overlayCanvas.width = width;
      overlayCanvas.height = height;

      // Draw video frame to canvas to get image data
      const tempCanvas = document.createElement("canvas");
      tempCanvas.width = width;
      tempCanvas.height = height;
      const tempCtx = tempCanvas.getContext("2d");
      if (!tempCtx) {
        animationFrameId = requestAnimationFrame(detectFrame);
        return;
      }
      tempCtx.drawImage(video, 0, 0);
      const imageData = tempCtx.getImageData(0, 0, width, height);

      // Process with OpenCV
      const src = cv.matFromImageData(imageData);
      const gray = new cv.Mat();
      const blurred = new cv.Mat();
      const edges = new cv.Mat();
      const contours = new cv.MatVector();
      const hierarchy = new cv.Mat();

      // Convert to grayscale
      cv.cvtColor(src, gray, cv.COLOR_RGBA2GRAY);
      
      // Apply Gaussian blur
      cv.GaussianBlur(gray, blurred, new cv.Size(5, 5), 0);
      
      // Detect edges
      cv.Canny(blurred, edges, 75, 200);
      
      // Find contours
      cv.findContours(edges, contours, hierarchy, cv.RETR_EXTERNAL, cv.CHAIN_APPROX_SIMPLE);

      // Find largest quadrilateral
      let maxArea = 0;
      let bestContour: Mat | null = null;

      for (let i = 0; i < contours.size(); i++) {
        const contour = contours.get(i);
        const area = cv.contourArea(contour);
        
        if (area > maxArea) {
          const peri = cv.arcLength(contour, true);
          const approx = new cv.Mat();
          cv.approxPolyDP(contour, approx, 0.02 * peri, true);
          
          if (approx.rows === 4 && area > (width * height * 0.1)) {
            maxArea = area;
            if (bestContour) bestContour.delete();
            bestContour = approx;
          } else {
            approx.delete();
          }
        }
      }

      // Clear overlay
      ctx.clearRect(0, 0, width, height);

      if (bestContour) {
        // Extract corner points
        const points: { x: number; y: number }[] = [];
        for (let i = 0; i < 4; i++) {
          const x = bestContour.data32F[i * 2];
          const y = bestContour.data32F[i * 2 + 1];
          points.push({ x, y });
        }
        
        const orderedPoints = orderPoints(points);
        detectedCorners.value = orderedPoints;

        // Draw polygon
        ctx.strokeStyle = "#00ff00";
        ctx.lineWidth = 3;
        ctx.beginPath();
        ctx.moveTo(orderedPoints[0].x, orderedPoints[0].y);
        for (let i = 1; i < 4; i++) {
          ctx.lineTo(orderedPoints[i].x, orderedPoints[i].y);
        }
        ctx.closePath();
        ctx.stroke();

        // Draw corner circles
        ctx.fillStyle = "#00ff00";
        orderedPoints.forEach(point => {
          ctx.beginPath();
          ctx.arc(point.x, point.y, 10, 0, 2 * Math.PI);
          ctx.fill();
        });

        bestContour.delete();
      } else {
        detectedCorners.value = [];
      }

      // Cleanup
      src.delete();
      gray.delete();
      blurred.delete();
      edges.delete();
      contours.delete();
      hierarchy.delete();

    } catch (e) {
      console.error("Edge detection error:", e);
    }

    animationFrameId = requestAnimationFrame(detectFrame);
  };

  detectFrame();
};

const stopEdgeDetection = () => {
  isDetecting.value = false;
  if (animationFrameId) {
    cancelAnimationFrame(animationFrameId);
    animationFrameId = null;
  }
};

const captureWithPerspectiveTransform = (): string | null => {
  if (!videoRef.value || detectedCorners.value.length !== 4 || !isOpenCVLoaded.value) {
    return null;
  }

  try {
    const video = videoRef.value;
    const width = video.videoWidth;
    const height = video.videoHeight;

    // Draw video to temp canvas
    const tempCanvas = document.createElement("canvas");
    tempCanvas.width = width;
    tempCanvas.height = height;
    const tempCtx = tempCanvas.getContext("2d");
    if (!tempCtx) return null;
    
    tempCtx.drawImage(video, 0, 0);
    const imageData = tempCtx.getImageData(0, 0, width, height);

    const src = cv.matFromImageData(imageData);
    const corners = detectedCorners.value;

    // Calculate output dimensions
    const widthA = Math.sqrt(Math.pow(corners[1].x - corners[0].x, 2) + Math.pow(corners[1].y - corners[0].y, 2));
    const widthB = Math.sqrt(Math.pow(corners[2].x - corners[3].x, 2) + Math.pow(corners[2].y - corners[3].y, 2));
    const maxWidth = Math.max(widthA, widthB);

    const heightA = Math.sqrt(Math.pow(corners[3].x - corners[0].x, 2) + Math.pow(corners[3].y - corners[0].y, 2));
    const heightB = Math.sqrt(Math.pow(corners[2].x - corners[1].x, 2) + Math.pow(corners[2].y - corners[1].y, 2));
    const maxHeight = Math.max(heightA, heightB);

    // Source points
    const srcPoints = cv.matFromArray(4, 1, cv.CV_32FC2, [
      corners[0].x, corners[0].y,
      corners[1].x, corners[1].y,
      corners[2].x, corners[2].y,
      corners[3].x, corners[3].y
    ]);

    // Destination points
    const dstPoints = cv.matFromArray(4, 1, cv.CV_32FC2, [
      0, 0,
      maxWidth - 1, 0,
      maxWidth - 1, maxHeight - 1,
      0, maxHeight - 1
    ]);

    // Get perspective transform
    const M = cv.getPerspectiveTransform(srcPoints, dstPoints);
    const dst = new cv.Mat();
    cv.warpPerspective(src, dst, M, new cv.Size(maxWidth, maxHeight));

    // Convert result to canvas
    const resultCanvas = document.createElement("canvas");
    resultCanvas.width = maxWidth;
    resultCanvas.height = maxHeight;
    const resultCtx = resultCanvas.getContext("2d");
    
    if (resultCtx) {
      const resultImageData = new ImageData(
        new Uint8ClampedArray(dst.data),
        maxWidth,
        maxHeight
      );
      resultCtx.putImageData(resultImageData, 0, 0);
    }

    // Cleanup
    src.delete();
    srcPoints.delete();
    dstPoints.delete();
    M.delete();
    dst.delete();

    return resultCanvas.toDataURL("image/jpeg", 0.85);
  } catch (e) {
    console.error("Perspective transform error:", e);
    return null;
  }
};

const captureImage = () => {
  if (!videoRef.value || !canvasRef.value || !isVideoReady.value) {
    error.value = "Kamera belum siap, tunggu sebentar...";
    return;
  }

  // Try perspective transform first if document detected
  if (detectedCorners.value.length === 4 && isOpenCVLoaded.value) {
    const transformedImage = captureWithPerspectiveTransform();
    if (transformedImage) {
      capturedImages.value.push(transformedImage);
      return;
    }
  }

  // Fallback to regular capture
  const video = videoRef.value;
  const canvas = canvasRef.value;
  canvas.width = video.videoWidth;
  canvas.height = video.videoHeight;

  const ctx = canvas.getContext("2d");
  if (ctx) {
    ctx.drawImage(video, 0, 0);
    const imageData = canvas.toDataURL("image/jpeg", 0.85);
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

watch(isVideoReady, (ready) => {
  if (ready && isOpenCVLoaded.value) {
    startEdgeDetection();
  }
});

watch(isOpenCVLoaded, (loaded) => {
  if (loaded && isVideoReady.value) {
    startEdgeDetection();
  }
});

onMounted(async () => {
  if (isMobile.value) {
    initScanner();
    startCamera();
  }
});

onUnmounted(() => {
  stopEdgeDetection();
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

    <div v-if="!isOpenCVLoaded" class="bg-yellow-500 text-white p-2 text-center text-sm">
      ⏳ Memuat scanner...
    </div>

    <div class="flex-1 relative overflow-hidden">
      <video
        ref="videoRef"
        autoplay
        playsinline
        class="w-full h-full object-cover"
      ></video>

      <!-- Overlay canvas for edge detection -->
      <canvas
        ref="overlayCanvasRef"
        class="absolute inset-0 w-full h-full pointer-events-none"
        style="object-fit: cover"
      ></canvas>

      <canvas ref="canvasRef" class="hidden"></canvas>

      <!-- Guide when no document detected -->
      <div
        v-if="isOpenCVLoaded && isVideoReady && detectedCorners.length < 4"
        class="absolute inset-8 border-2 border-dashed border-white/50 rounded-lg pointer-events-none"
      >
        <div class="absolute inset-0 flex items-center justify-center">
          <p class="bg-black/50 text-white px-3 py-1 rounded text-sm">
            Arahkan kamera ke dokumen
          </p>
        </div>
      </div>

      <!-- Document detected -->
      <div
        v-if="detectedCorners.length === 4"
        class="absolute top-4 left-1/2 -translate-x-1/2 bg-green-500 text-white px-4 py-2 rounded-full text-sm font-medium shadow-lg"
      >
        ✓ Dokumen terdeteksi
      </div>

      <div class="absolute bottom-4 left-0 right-0 flex justify-center gap-4">
        <button
          @click="captureImage"
          :disabled="!isVideoReady"
          class="w-16 h-16 bg-white rounded-full border-4 flex items-center justify-center disabled:opacity-50 transition-all active:scale-95"
          :class="detectedCorners.length === 4 ? 'border-green-400' : 'border-gray-300'"
        >
          <div
            class="w-12 h-12 rounded-full transition-colors"
            :class="detectedCorners.length === 4 ? 'bg-green-500' : 'bg-red-500'"
          ></div>
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
      <div class="flex gap-2 mt-3 items-center">
        <span class="text-white text-sm">{{ capturedImages.length }} halaman</span>
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
