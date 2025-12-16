<script setup lang="ts">
import { computed } from "vue";

const props = defineProps<{
  modelValue: number;
  label?: string;
  error?: string;
  required?: boolean;
  disabled?: boolean;
}>();

const emit = defineEmits<{
  "update:modelValue": [value: number];
}>();

const formatCurrency = (value: number): string => {
  return new Intl.NumberFormat("id-ID").format(value);
};

const displayValue = computed(() => formatCurrency(props.modelValue || 0));

const handleInput = (event: Event) => {
  const input = event.target as HTMLInputElement;
  const rawValue = input.value.replace(/\D/g, "");
  const numValue = parseInt(rawValue) || 0;
  emit("update:modelValue", numValue);
};
</script>

<template>
  <div class="mb-4">
    <label v-if="label" class="block text-sm font-medium text-gray-700 mb-1">
      {{ label }}
      <span v-if="required" class="text-red-500">*</span>
    </label>
    <div class="relative">
      <span class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-500"
        >Rp</span
      >
      <input
        type="text"
        :value="displayValue"
        :disabled="disabled"
        @input="handleInput"
        class="w-full pl-10 pr-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all disabled:bg-gray-100"
        :class="error ? 'border-red-500' : 'border-gray-300'"
      />
    </div>
    <p v-if="error" class="mt-1 text-sm text-red-500">{{ error }}</p>
  </div>
</template>
