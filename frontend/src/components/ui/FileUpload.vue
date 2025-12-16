<script setup lang="ts">
import { ref } from "vue";

const props = defineProps<{
  label?: string;
  accept?: string;
  error?: string;
  required?: boolean;
  disabled?: boolean;
}>();

const emit = defineEmits<{
  "update:file": [file: File | null];
}>();

const isDragging = ref(false);
const selectedFile = ref<File | null>(null);
const previewUrl = ref<string | null>(null);
const inputId = ref(
  `file-upload-${Math.random().toString(36).substring(2, 9)}`
);

const handleFile = (file: File) => {
  selectedFile.value = file;
  emit("update:file", file);

  if (file.type.startsWith("image/")) {
    previewUrl.value = URL.createObjectURL(file);
  } else {
    previewUrl.value = null;
  }
};

const onDrop = (e: DragEvent) => {
  isDragging.value = false;
  const file = e.dataTransfer?.files[0];
  if (file) handleFile(file);
};

const onFileSelect = (e: Event) => {
  const input = e.target as HTMLInputElement;
  const file = input.files?.[0];
  if (file) handleFile(file);
};

const clearFile = () => {
  selectedFile.value = null;
  previewUrl.value = null;
  emit("update:file", null);
};
</script>

<template>
  <div class="mb-4">
    <label v-if="label" class="block text-sm font-medium text-gray-700 mb-1">
      {{ label }}
      <span v-if="required" class="text-red-500">*</span>
    </label>

    <div
      v-if="!selectedFile"
      @dragover.prevent="isDragging = true"
      @dragleave="isDragging = false"
      @drop.prevent="onDrop"
      class="border-2 border-dashed rounded-lg p-6 text-center cursor-pointer transition-colors"
      :class="[
        isDragging ? 'border-blue-500 bg-blue-50' : 'border-gray-300',
        error ? 'border-red-500' : '',
        disabled ? 'bg-gray-100 cursor-not-allowed' : 'hover:border-blue-400',
      ]"
    >
      <input
        type="file"
        :accept="accept"
        :disabled="disabled"
        @change="onFileSelect"
        class="hidden"
        :id="inputId"
      />
      <label :for="inputId" class="cursor-pointer">
        <div class="text-4xl mb-2">📁</div>
        <p class="text-gray-600">Drag & drop file atau klik untuk memilih</p>
        <p v-if="accept" class="text-sm text-gray-400 mt-1">{{ accept }}</p>
      </label>
    </div>

    <div v-else class="border rounded-lg p-4 flex items-center justify-between">
      <div class="flex items-center gap-3">
        <img
          v-if="previewUrl"
          :src="previewUrl"
          class="w-12 h-12 object-cover rounded"
        />
        <div
          v-else
          class="w-12 h-12 bg-gray-100 rounded flex items-center justify-center text-2xl"
        >
          📄
        </div>
        <div>
          <p class="font-medium text-gray-700">{{ selectedFile.name }}</p>
          <p class="text-sm text-gray-400">
            {{ (selectedFile.size / 1024).toFixed(1) }} KB
          </p>
        </div>
      </div>
      <button
        @click="clearFile"
        type="button"
        class="text-red-500 hover:text-red-700"
      >
        ✕
      </button>
    </div>

    <p v-if="error" class="mt-1 text-sm text-red-500">{{ error }}</p>
  </div>
</template>
