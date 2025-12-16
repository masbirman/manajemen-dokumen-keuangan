import { ref } from "vue";

interface Toast {
  id: number;
  message: string;
  type: "success" | "error" | "warning" | "info";
}

const toasts = ref<Toast[]>([]);
let nextId = 0;

export function useToast() {
  const addToast = (
    message: string,
    type: "success" | "error" | "warning" | "info" = "info"
  ) => {
    const id = nextId++;
    toasts.value.push({ id, message, type });
  };

  const removeToast = (id: number) => {
    toasts.value = toasts.value.filter((t) => t.id !== id);
  };

  const success = (message: string) => addToast(message, "success");
  const error = (message: string) => addToast(message, "error");
  const warning = (message: string) => addToast(message, "warning");
  const info = (message: string) => addToast(message, "info");

  return {
    toasts,
    addToast,
    removeToast,
    success,
    error,
    warning,
    info,
  };
}
