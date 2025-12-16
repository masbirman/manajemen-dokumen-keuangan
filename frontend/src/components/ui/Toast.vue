<script setup lang="ts">
import { ref, onMounted } from "vue";

const props = defineProps<{
  message: string;
  type?: "success" | "error" | "warning" | "info";
  duration?: number;
}>();

const emit = defineEmits<{
  close: [];
}>();

const visible = ref(true);

const typeStyles = {
  success: "bg-green-500",
  error: "bg-red-500",
  warning: "bg-yellow-500",
  info: "bg-blue-500",
};

const icons = {
  success: "✓",
  error: "✕",
  warning: "⚠",
  info: "ℹ",
};

onMounted(() => {
  setTimeout(() => {
    visible.value = false;
    emit("close");
  }, props.duration || 3000);
});
</script>

<template>
  <Transition name="toast">
    <div
      v-if="visible"
      class="fixed top-4 right-4 z-50 flex items-center gap-3 px-4 py-3 rounded-lg text-white shadow-lg"
      :class="typeStyles[type || 'info']"
    >
      <span class="text-lg">{{ icons[type || "info"] }}</span>
      <span>{{ message }}</span>
      <button
        @click="
          visible = false;
          $emit('close');
        "
        class="ml-2 hover:opacity-80"
      >
        ✕
      </button>
    </div>
  </Transition>
</template>

<style scoped>
.toast-enter-active,
.toast-leave-active {
  transition: all 0.3s ease;
}
.toast-enter-from {
  opacity: 0;
  transform: translateX(100%);
}
.toast-leave-to {
  opacity: 0;
  transform: translateY(-100%);
}
</style>
