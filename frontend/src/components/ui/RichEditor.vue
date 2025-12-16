<script setup lang="ts">
import { QuillEditor } from "@vueup/vue-quill";
import "@vueup/vue-quill/dist/vue-quill.snow.css";

defineProps<{
  modelValue: string;
  label?: string;
  placeholder?: string;
  error?: string;
  required?: boolean;
}>();

const emit = defineEmits<{
  "update:modelValue": [value: string];
}>();

const toolbarOptions = [
  ["bold", "italic", "underline", "strike"],
  [{ list: "ordered" }, { list: "bullet" }],
  [{ header: [1, 2, 3, false] }],
  [{ align: [] }], // alignment: left, center, right, justify
  [{ color: [] }, { background: [] }],
  ["link", "image"],
  ["clean"],
];

const handleUpdate = (content: string) => {
  emit("update:modelValue", content);
};
</script>

<template>
  <div class="mb-4">
    <label v-if="label" class="block text-sm font-medium text-gray-700 mb-1">
      {{ label }}
      <span v-if="required" class="text-red-500">*</span>
    </label>
    <div
      class="border rounded-lg overflow-hidden"
      :class="error ? 'border-red-500' : 'border-gray-300'"
    >
      <QuillEditor
        theme="snow"
        :content="modelValue"
        content-type="html"
        :toolbar="toolbarOptions"
        :placeholder="placeholder || 'Tulis konten di sini...'"
        @update:content="handleUpdate"
        style="min-height: 200px"
      />
    </div>
    <p v-if="error" class="mt-1 text-sm text-red-500">{{ error }}</p>
  </div>
</template>

<style>
.ql-container {
  font-size: 14px;
  min-height: 150px;
}
.ql-editor {
  min-height: 150px;
}
</style>
