<template>
  <div class="loadable-select">
    <div class="select-wrapper" @click="toggleDropdown">
      <div class="select-display">
        <span v-if="modelValue" class="selected-value">
          {{ getDisplayValue(modelValue) }}
        </span>
        <span v-else class="placeholder">{{ placeholder }}</span>
      </div>
      <div class="select-arrow" :class="{ 'rotated': isOpen }">
        <svg viewBox="0 0 24 24" fill="currentColor">
          <path d="M7 10l5 5 5-5z" />
        </svg>
      </div>
    </div>

    <div class="dropdown" v-show="isOpen">
      <!-- 搜索框 -->
      <div class="search-container">
        <div class="search-input-wrapper">
          <svg class="search-icon" viewBox="0 0 24 24" fill="currentColor">
            <path d="M15.5 14h-.79l-.28-.27C15.41 12.59 16 11.11 16 9.5 16 5.91 13.09 3 9.5 3S3 5.91 3 9.5 5.91 16 9.5 16c1.61 0 3.09-.59 4.23-1.57l.27.28v.79l5 4.99L20.49 19l-4.99-5zm-6 0C7.01 14 5 11.99 5 9.5S7.01 5 9.5 5 14 7.01 14 9.5 11.99 14 9.5 14z"/>
          </svg>
          <input
            ref="searchInput"
            v-model="searchQuery"
            type="text"
            class="search-input"
            placeholder="搜索音色..."
            @click.stop
            @input="handleSearch"
            @keydown="handleKeydown"
          />
          <button
            v-if="searchQuery"
            class="clear-search-btn"
            @click.stop="clearSearch"
            title="清除搜索"
          >
            <svg viewBox="0 0 24 24" fill="currentColor">
              <path d="M19 6.41L17.59 5 12 10.59 6.41 5 5 6.41 10.59 12 5 17.59 6.41 19 12 13.41 17.59 19 19 17.59 13.41 12z"/>
            </svg>
          </button>
        </div>
      </div>

      <div class="dropdown-content">
        <!-- 搜索结果统计 -->
        <div v-if="searchQuery && !loading" class="search-stats">
          <span>{{ filteredItems.length }} 个结果</span>
        </div>

        <div
          v-for="item in displayItems"
          :key="getItemKey(item)"
          class="dropdown-item"
          :class="{ 'selected': isItemSelected(item) }"
          @click="selectItem(item)"
        >
          <span v-if="searchQuery" v-html="highlightText(getItemLabel(item), searchQuery)"></span>
          <span v-else>{{ getItemLabel(item) }}</span>
        </div>

        <!-- 无搜索结果 -->
        <div v-if="searchQuery && filteredItems.length === 0 && !loading" class="no-results">
          未找到匹配的音色
        </div>

        <!-- 加载更多按钮 -->
        <div
          v-if="hasMore && !loading && !searchQuery"
          class="load-more-btn"
          @click.stop="loadMore"
        >
          加载更多
        </div>

        <!-- 加载状态 -->
        <div v-if="loading" class="loading-state" @click.stop>
          <div class="loading-spinner"></div>
          <span>加载中...</span>
        </div>

        <!-- 无更多数据 -->
        <div v-if="!hasMore && displayItems.length > 0 && !searchQuery" class="no-more" @click.stop>
          没有更多数据
        </div>

        <!-- 空状态 -->
        <div v-if="!loading && displayItems.length === 0 && !searchQuery" class="empty-state" @click.stop>
          {{ emptyText }}
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted, onUnmounted, nextTick } from 'vue';

const props = defineProps({
  modelValue: [String, Number],
  items: {
    type: Array,
    default: () => []
  },
  placeholder: {
    type: String,
    default: '请选择'
  },
  emptyText: {
    type: String,
    default: '暂无数据'
  },
  valueKey: {
    type: String,
    default: 'value'
  },
  labelKey: {
    type: String,
    default: 'label'
  },
  hasMore: {
    type: Boolean,
    default: false
  },
  loading: {
    type: Boolean,
    default: false
  }
});

const emit = defineEmits(['update:modelValue', 'load-more']);

const isOpen = ref(false);
const searchQuery = ref('');
const searchInput = ref(null);

// 计算属性 - 过滤后的项目
const filteredItems = computed(() => {
  if (!searchQuery.value.trim()) {
    return props.items;
  }

  const query = searchQuery.value.toLowerCase().trim();
  return props.items.filter(item => {
    const label = getItemLabel(item).toLowerCase();
    return label.includes(query);
  });
});

// 高亮匹配的文本
const highlightText = (text, query) => {
  if (!query.trim()) return text;

  const regex = new RegExp(`(${query.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')})`, 'gi');
  return text.replace(regex, '<mark>$1</mark>');
};

// 计算属性 - 显示的项目（用于渲染）
const displayItems = computed(() => filteredItems.value);

// 获取显示值
const getDisplayValue = (value) => {
  const item = props.items.find(item => getItemValue(item) === value);
  return item ? getItemLabel(item) : value;
};

// 获取项目值
const getItemValue = (item) => {
  return typeof item === 'object' ? item[props.valueKey] : item;
};

// 获取项目标签
const getItemLabel = (item) => {
  return typeof item === 'object' ? item[props.labelKey] : item;
};

// 获取项目键
const getItemKey = (item) => {
  return getItemValue(item);
};

// 检查项目是否被选中
const isItemSelected = (item) => {
  return getItemValue(item) === props.modelValue;
};

// 切换下拉框
const toggleDropdown = () => {
  isOpen.value = !isOpen.value;
  if (isOpen.value) {
    // 打开下拉框时，延迟聚焦搜索框
    nextTick(() => {
      if (searchInput.value) {
        searchInput.value.focus();
      }
    });
  } else {
    // 关闭下拉框时，清空搜索
    clearSearch();
  }
};

// 选择项目
const selectItem = (item) => {
  const value = getItemValue(item);
  emit('update:modelValue', value);
  isOpen.value = false;
  clearSearch();
};

// 处理搜索
const handleSearch = () => {
  // 搜索逻辑已在 computed 中处理
};

// 处理键盘事件
const handleKeydown = (event) => {
  if (isOpen.value) {
    if (event.key === 'Escape') {
      isOpen.value = false;
      clearSearch();
    } else if (event.key === 'Enter' && searchInput.value) {
      // 如果搜索框有焦点，则触发搜索
      searchInput.value.focus();
    }
  }
};

// 清除搜索
const clearSearch = () => {
  searchQuery.value = '';
  // 重新聚焦搜索框
  nextTick(() => {
    if (searchInput.value) {
      searchInput.value.focus();
    }
  });
};

// 加载更多
const loadMore = () => {
  if (!props.loading && props.hasMore) {
    emit('load-more');
    // 不关闭下拉框，让用户可以继续浏览
  }
};

// 点击外部关闭下拉框
const handleClickOutside = (event) => {
  if (!event.target.closest('.loadable-select')) {
    isOpen.value = false;
    clearSearch();
  }
};

// 监听点击事件
onMounted(() => {
  document.addEventListener('click', handleClickOutside);
});

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside);
});
</script>

<style scoped>
.loadable-select {
  position: relative;
  width: 100%;
}

.select-wrapper {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 12px;
  border: 1px solid rgba(255, 255, 255, 0.3);
  border-radius: 6px;
  background: rgba(255, 255, 255, 0.1);
  color: white;
  font-size: 14px;
  cursor: pointer;
  transition: border-color 0.2s ease;
}

.select-wrapper:hover {
  border-color: rgba(255, 255, 255, 0.6);
}

.select-display {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.selected-value {
  color: white;
}

.placeholder {
  color: rgba(255, 255, 255, 0.5);
}

.select-arrow {
  margin-left: 8px;
  transition: transform 0.2s ease;
}

.select-arrow.rotated {
  transform: rotate(180deg);
}

.select-arrow svg {
  width: 16px;
  height: 16px;
}

.dropdown {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(0, 0, 0, 0.1);
  border-radius: 6px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  z-index: 1000;
  margin-top: 4px;
  max-height: 200px;
  overflow-y: auto;
}

.dropdown-content {
  padding: 4px 0;
}

.dropdown-item {
  padding: 8px 12px;
  cursor: pointer;
  transition: background-color 0.2s ease;
  color: #333;
  font-size: 14px;
}

.dropdown-item:hover {
  background-color: rgba(24, 144, 255, 0.1);
}

.dropdown-item.selected {
  background-color: rgba(24, 144, 255, 0.2);
  color: #1890ff;
  font-weight: 500;
}

.dropdown-item mark {
  background-color: #ffeb3b;
  color: #333;
  padding: 0 1px;
  border-radius: 2px;
  font-weight: bold;
}

.load-more-btn {
  padding: 8px 12px;
  text-align: center;
  color: #1890ff;
  cursor: pointer;
  font-size: 14px;
  border-top: 1px solid rgba(0, 0, 0, 0.1);
  transition: background-color 0.2s ease;
}

.load-more-btn:hover {
  background-color: rgba(24, 144, 255, 0.1);
}

.loading-state {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 12px;
  color: #666;
  font-size: 14px;
}

.loading-spinner {
  width: 16px;
  height: 16px;
  border: 2px solid rgba(0, 0, 0, 0.1);
  border-top: 2px solid #1890ff;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

.no-more {
  padding: 8px 12px;
  text-align: center;
  color: #999;
  font-size: 12px;
  border-top: 1px solid rgba(0, 0, 0, 0.1);
}

.empty-state {
  padding: 12px;
  text-align: center;
  color: #999;
  font-size: 14px;
}

.search-container {
  padding: 8px 12px;
  border-bottom: 1px solid rgba(0, 0, 0, 0.1);
  background: rgba(248, 249, 250, 0.95);
}

.search-input-wrapper {
  position: relative;
  display: flex;
  align-items: center;
  background-color: white;
  border: 1px solid rgba(0, 0, 0, 0.1);
  border-radius: 6px;
  padding: 6px 10px;
  transition: border-color 0.2s ease;
}

.search-input-wrapper:focus-within {
  border-color: #1890ff;
  box-shadow: 0 0 0 2px rgba(24, 144, 255, 0.1);
}

.search-icon {
  width: 16px;
  height: 16px;
  margin-right: 8px;
  color: #999;
  flex-shrink: 0;
}

.search-input {
  flex: 1;
  border: none;
  outline: none;
  background-color: transparent;
  font-size: 14px;
  color: #333;
  min-width: 0;
}

.search-input::placeholder {
  color: #999;
}

.clear-search-btn {
  background: none;
  border: none;
  cursor: pointer;
  padding: 2px;
  margin-left: 4px;
  color: #999;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s ease;
  flex-shrink: 0;
}

.clear-search-btn:hover {
  background-color: rgba(0, 0, 0, 0.1);
  color: #666;
}

.clear-search-btn svg {
  width: 14px;
  height: 14px;
}

.search-stats {
  padding: 6px 12px;
  font-size: 12px;
  color: #666;
  background: rgba(248, 249, 250, 0.5);
  border-bottom: 1px solid rgba(0, 0, 0, 0.05);
}

.no-results {
  padding: 12px;
  text-align: center;
  color: #999;
  font-size: 14px;
  font-style: italic;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

/* 滚动条样式 */
.dropdown::-webkit-scrollbar {
  width: 6px;
}

.dropdown::-webkit-scrollbar-track {
  background: rgba(0, 0, 0, 0.1);
  border-radius: 3px;
}

.dropdown::-webkit-scrollbar-thumb {
  background: rgba(0, 0, 0, 0.3);
  border-radius: 3px;
}

.dropdown::-webkit-scrollbar-thumb:hover {
  background: rgba(0, 0, 0, 0.5);
}
</style>
