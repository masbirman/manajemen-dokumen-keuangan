<script setup lang="ts">
import { computed, ref, watch } from "vue";

interface Column {
  key: string;
  label: string;
  sortable?: boolean;
  width?: string;
}

const props = defineProps<{
  columns: Column[];
  data: Record<string, unknown>[];
  loading?: boolean;
  sortKey?: string;
  sortOrder?: "asc" | "desc";
  currentPage?: number;
  totalPages?: number;
  totalItems?: number;
  perPage?: number;
  selectable?: boolean;
  selectedIds?: (string | number)[];
  idKey?: string;
  searchable?: boolean;
  searchPlaceholder?: string;
}>();

const emit = defineEmits<{
  sort: [key: string];
  "page-change": [page: number];
  "per-page-change": [perPage: number];
  "row-click": [row: Record<string, unknown>];
  "selection-change": [ids: (string | number)[]];
  search: [query: string];
}>();

const searchQuery = ref("");
const perPageOptions = [10, 25, 50, 100];

const handleSearch = () => {
  emit("search", searchQuery.value);
};

const handlePerPageChange = (e: Event) => {
  const value = parseInt((e.target as HTMLSelectElement).value);
  emit("per-page-change", value);
};

const localSelectedIds = ref<(string | number)[]>([]);

watch(
  () => props.selectedIds,
  (newVal) => {
    localSelectedIds.value = newVal || [];
  },
  { immediate: true }
);

const idField = computed(() => props.idKey || "id");

const allSelected = computed(() => {
  if (props.data.length === 0) return false;
  return props.data.every((row) =>
    localSelectedIds.value.includes(row[idField.value] as string | number)
  );
});

const someSelected = computed(() => {
  return (
    localSelectedIds.value.length > 0 &&
    props.data.some((row) =>
      localSelectedIds.value.includes(row[idField.value] as string | number)
    )
  );
});

const toggleAll = () => {
  if (allSelected.value) {
    // Deselect all current page items
    const currentIds = props.data.map(
      (row) => row[idField.value] as string | number
    );
    localSelectedIds.value = localSelectedIds.value.filter(
      (id) => !currentIds.includes(id)
    );
  } else {
    // Select all current page items
    const currentIds = props.data.map(
      (row) => row[idField.value] as string | number
    );
    const newIds = [...localSelectedIds.value];
    currentIds.forEach((id) => {
      if (!newIds.includes(id)) newIds.push(id);
    });
    localSelectedIds.value = newIds;
  }
  emit("selection-change", localSelectedIds.value);
};

const toggleRow = (row: Record<string, unknown>) => {
  const id = row[idField.value] as string | number;
  const index = localSelectedIds.value.indexOf(id);
  if (index > -1) {
    localSelectedIds.value.splice(index, 1);
  } else {
    localSelectedIds.value.push(id);
  }
  emit("selection-change", [...localSelectedIds.value]);
};

const isSelected = (row: Record<string, unknown>) => {
  return localSelectedIds.value.includes(row[idField.value] as string | number);
};

const handleSort = (column: Column) => {
  if (column.sortable) {
    emit("sort", column.key);
  }
};

const pages = computed(() => {
  const total = props.totalPages || 1;
  const current = props.currentPage || 1;
  const pages: (number | string)[] = [];

  if (total <= 7) {
    for (let i = 1; i <= total; i++) pages.push(i);
  } else {
    if (current <= 3) {
      pages.push(1, 2, 3, 4, "...", total);
    } else if (current >= total - 2) {
      pages.push(1, "...", total - 3, total - 2, total - 1, total);
    } else {
      pages.push(1, "...", current - 1, current, current + 1, "...", total);
    }
  }
  return pages;
});
</script>

<template>
  <div
    class="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden"
  >
    <!-- Search and filters bar -->
    <div
      v-if="searchable || $slots.filters"
      class="px-4 py-3 border-b border-gray-200 flex flex-wrap gap-3 items-center"
    >
      <div v-if="searchable" class="flex-1 min-w-[200px] max-w-md">
        <div class="relative">
          <input
            v-model="searchQuery"
            type="text"
            :placeholder="searchPlaceholder || 'Cari...'"
            class="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg text-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            @keyup.enter="handleSearch"
          />
          <span class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400"
            >🔍</span
          >
        </div>
      </div>
      <slot name="filters"></slot>
      <button
        v-if="searchable"
        @click="handleSearch"
        class="px-4 py-2 bg-blue-600 text-white text-sm rounded-lg hover:bg-blue-700"
      >
        Cari
      </button>
    </div>

    <!-- Bulk action bar -->
    <div
      v-if="selectable && localSelectedIds.length > 0"
      class="px-4 py-2 bg-blue-50 border-b border-blue-200 flex items-center gap-4"
    >
      <span class="text-sm text-blue-700 font-medium">
        {{ localSelectedIds.length }} item dipilih
      </span>
      <slot name="bulk-actions" :selectedIds="localSelectedIds"></slot>
    </div>

    <div class="overflow-x-auto">
      <table class="w-full">
        <thead class="bg-gray-50 border-b border-gray-200">
          <tr>
            <th v-if="selectable" class="px-4 py-3 w-12">
              <input
                type="checkbox"
                :checked="allSelected"
                :indeterminate="someSelected && !allSelected"
                @change="toggleAll"
                class="w-4 h-4 rounded border-gray-300 text-blue-600 focus:ring-blue-500"
              />
            </th>
            <th
              v-for="column in columns"
              :key="column.key"
              @click="handleSort(column)"
              class="px-4 py-3 text-left text-sm font-semibold text-gray-700"
              :class="{ 'cursor-pointer hover:bg-gray-100': column.sortable }"
              :style="column.width ? { width: column.width } : {}"
            >
              <div class="flex items-center gap-1">
                {{ column.label }}
                <span
                  v-if="column.sortable && sortKey === column.key"
                  class="text-blue-600"
                >
                  {{ sortOrder === "asc" ? "↑" : "↓" }}
                </span>
              </div>
            </th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="loading">
            <td
              :colspan="selectable ? columns.length + 1 : columns.length"
              class="px-4 py-8 text-center text-gray-500"
            >
              <div class="flex items-center justify-center gap-2">
                <span class="animate-spin">⏳</span> Loading...
              </div>
            </td>
          </tr>
          <tr v-else-if="data.length === 0">
            <td
              :colspan="selectable ? columns.length + 1 : columns.length"
              class="px-4 py-8 text-center text-gray-500"
            >
              Tidak ada data
            </td>
          </tr>
          <tr
            v-else
            v-for="(row, index) in data"
            :key="index"
            @click="$emit('row-click', row)"
            class="border-b border-gray-100 hover:bg-gray-50 cursor-pointer transition-colors"
            :class="{ 'bg-blue-50': selectable && isSelected(row) }"
          >
            <td v-if="selectable" class="px-4 py-3" @click.stop>
              <input
                type="checkbox"
                :checked="isSelected(row)"
                @change="toggleRow(row)"
                class="w-4 h-4 rounded border-gray-300 text-blue-600 focus:ring-blue-500"
              />
            </td>
            <td
              v-for="column in columns"
              :key="column.key"
              class="px-4 py-3 text-sm"
            >
              <slot :name="column.key" :row="row" :value="row[column.key]">
                {{ row[column.key] }}
              </slot>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <div
      v-if="totalPages && totalPages >= 1"
      class="px-4 py-3 border-t border-gray-200 flex flex-wrap justify-between items-center gap-3"
    >
      <div class="flex items-center gap-3">
        <span class="text-sm text-gray-500"> Tampilkan </span>
        <select
          :value="perPage"
          @change="handlePerPageChange"
          class="border border-gray-300 rounded px-2 py-1 text-sm focus:ring-2 focus:ring-blue-500"
        >
          <option v-for="opt in perPageOptions" :key="opt" :value="opt">
            {{ opt }}
          </option>
        </select>
        <span class="text-sm text-gray-500"> data per halaman </span>
      </div>
      <span class="text-sm text-gray-500">
        Halaman {{ currentPage }} dari {{ totalPages }}
        <template v-if="totalItems"> ({{ totalItems }} data)</template>
      </span>
      <div class="flex gap-1">
        <button
          @click="currentPage > 1 && $emit('page-change', currentPage - 1)"
          :disabled="currentPage <= 1"
          class="px-3 py-1 text-sm rounded hover:bg-gray-100 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          ←
        </button>
        <button
          v-for="page in pages"
          :key="page"
          @click="typeof page === 'number' && $emit('page-change', page)"
          :disabled="page === '...'"
          class="px-3 py-1 text-sm rounded"
          :class="
            page === currentPage
              ? 'bg-blue-600 text-white'
              : page === '...'
              ? 'text-gray-400 cursor-default'
              : 'hover:bg-gray-100'
          "
        >
          {{ page }}
        </button>
        <button
          @click="
            currentPage < totalPages && $emit('page-change', currentPage + 1)
          "
          :disabled="currentPage >= totalPages"
          class="px-3 py-1 text-sm rounded hover:bg-gray-100 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          →
        </button>
      </div>
    </div>
  </div>
</template>
