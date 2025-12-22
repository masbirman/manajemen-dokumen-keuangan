import { defineStore } from 'pinia';
import { ref } from 'vue';

export interface Notification {
  id: number;
  message: string;
  type: 'info' | 'success' | 'error' | 'warning';
}

export const useNotificationStore = defineStore('notification', () => {
  const notifications = ref<Notification[]>([]);

  function show(message: string, type: 'info' | 'success' | 'error' | 'warning' = 'info', duration = 5000) {
    const id = Date.now();
    notifications.value.push({ id, message, type });
    
    if (duration > 0) {
      setTimeout(() => {
        remove(id);
      }, duration);
    }
  }

  function success(message: string, duration = 5000) {
    show(message, 'success', duration);
  }

  function error(message: string, duration = 5000) {
    show(message, 'error', duration);
  }

  function warning(message: string, duration = 5000) {
    show(message, 'warning', duration);
  }

  function info(message: string, duration = 5000) {
    show(message, 'info', duration);
  }

  function remove(id: number) {
    const index = notifications.value.findIndex(n => n.id === id);
    if (index !== -1) {
      notifications.value.splice(index, 1);
    }
  }

  return {
    notifications,
    show,
    success,
    error,
    warning,
    info,
    remove
  };
});
