<script setup lang="ts">
import { ref, computed, watch } from "vue";

interface Option {
  value: string | number;
  label: string;
}

const props = defineProps<{
  modelValue: string | number | null;
  options: Option[];
  label?: string;
  placeholder?: string;
  error?: string;
  required?: boolean;
  disabled?: boolean;
  searchable?: boolean;
}>();

const emit = defineEmits<{
  "update:modelValue": [value: string | number | null];
}>();

const search = ref("");
const isOpen = ref(false);

const filteredOptions = computed(() => {
  if (!props.searchable || !search.value) return props.options;
  return props.options.filter((opt) =>
    opt.label.toLowerCase().includes(search.value.toLowerCase())
  );
});

const selectedLabel = computed(() => {
  const selected = props.options.find((opt) => opt.value === props.modelValue);
  return selected?.label || "";
});

const selectOption = (value: string | number) => {
  emit("update:modelValue", value);
  isOpen.value = false;
  search.value = "";
};

watch(isOpen, (val) => {
  if (!val) search.value = "";
});
</script>

<template>
  <div class="mb-4 relative">
    <label v-if="label" class="block text-sm font-medium text-gray-700 mb-1">
      {{ label }}
      <span v-if="required" class="text-red-500">*</span>
    </label>
    <div
      @click="!disabled && (isOpen = !isOpen)"
      class="w-full px-3 py-2 border rounded-lg cursor-pointer flex justify-between items-center"
      :class="[
        error ? 'border-red-500' : 'border-gray-300',
        disabled ? 'bg-gray-100 cursor-not-allowed' : 'bg-white',
      ]"
    >
      <span :class="selectedLabel ? 'text-gray-900' : 'text-gray-400'">
        {{ selectedLabel || placeholder || "Pilih..." }}
      </span>
      <span class="text-gray-400">▼</span>
    </div>

    <div
      v-if="isOpen"
      class="absolute z-50 w-full mt-1 bg-white border border-gray-300 rounded-lg shadow-lg max-h-60 overflow-auto"
    >
      <input
        v-if="searchable"
        v-model="search"
        type="text"
        placeholder="Cari..."
        class="w-full px-3 py-2 border-b border-gray-200 focus:outline-none"
        @click.stop
      />
      <div
        v-for="option in filteredOptions"
        :key="option.value"
        @click="selectOption(option.value)"
        class="px-3 py-2 hover:bg-blue-50 cursor-pointer"
        :class="{ 'bg-blue-100': option.value === modelValue }"
      >
        {{ option.label }}
      </div>
      <div v-if="filteredOptions.length === 0" class="px-3 py-2 text-gray-400">
        Tidak ada data
      </div>
    </div>

    <p v-if="error" class="mt-1 text-sm text-red-500">{{ error }}</p>
  </div>
</template>
