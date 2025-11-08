# shadcn-vue 集成指南

本项目已经集成了 shadcn/ui 的 Vue 版本（shadcn-vue），提供现代化、美观的 UI 组件。

## ✅ 已完成的配置

- ✅ 安装了 Tailwind CSS
- ✅ 安装了 shadcn-vue 相关依赖
- ✅ 配置了路径别名 `@/`
- ✅ 创建了工具函数 `cn()`
- ✅ 添加了全局样式和 CSS 变量
- ✅ 创建了 `components.json` 配置文件

## 🎨 添加组件

shadcn-vue 提供了大量预制组件。你可以通过以下方式添加：

### 方法 1：使用 shadcn-vue CLI（推荐）

首先安装 CLI：

```bash
npm install -g shadcn-vue
```

然后添加组件：

```bash
# 在 frontend 目录下执行
npx shadcn-vue@latest add button
npx shadcn-vue@latest add card
npx shadcn-vue@latest add input
npx shadcn-vue@latest add dialog
npx shadcn-vue@latest add table
```

### 方法 2：手动复制组件

访问 [shadcn-vue 组件库](https://www.shadcn-vue.com/docs/components/accordion) 并手动复制组件代码到 `src/components/ui/` 目录。

## 常用组件示例

### Button 按钮

```vue
<script setup>
import { Button } from '@/components/ui/button'
</script>

<template>
  <Button>Click me</Button>
  <Button variant="outline">Outline</Button>
  <Button variant="destructive">Delete</Button>
</template>
```

### Card 卡片

```vue
<script setup>
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'
</script>

<template>
  <Card>
    <CardHeader>
      <CardTitle>Card Title</CardTitle>
      <CardDescription>Card description</CardDescription>
    </CardHeader>
    <CardContent>
      <p>Card content goes here</p>
    </CardContent>
    <CardFooter>
      <Button>Footer action</Button>
    </CardFooter>
  </Card>
</template>
```

### Input 输入框

```vue
<script setup>
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
</script>

<template>
  <div class="space-y-2">
    <Label for="email">Email</Label>
    <Input id="email" type="email" placeholder="Enter your email" />
  </div>
</template>
```

### Table 表格

```vue
<script setup>
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
</script>

<template>
  <Table>
    <TableHeader>
      <TableRow>
        <TableHead>Name</TableHead>
        <TableHead>Value</TableHead>
      </TableRow>
    </TableHeader>
    <TableBody>
      <TableRow>
        <TableCell>ETH</TableCell>
        <TableCell>$3,000</TableCell>
      </TableRow>
    </TableBody>
  </Table>
</template>
```

## 🎭 主题切换

shadcn-vue 支持亮色/暗色主题。切换主题只需在根元素添加/移除 `dark` 类：

```javascript
// 切换到暗色模式
document.documentElement.classList.add('dark')

// 切换到亮色模式
document.documentElement.classList.remove('dark')
```

## 🎨 自定义样式

所有颜色和样式变量在 `src/assets/index.css` 中定义。你可以修改这些 CSS 变量来自定义主题：

```css
:root {
  --primary: 222.2 47.4% 11.2%;  /* 主色调 */
  --secondary: 210 40% 96.1%;     /* 次要色调 */
  --destructive: 0 84.2% 60.2%;   /* 危险色（删除等） */
  --radius: 0.5rem;                /* 圆角大小 */
}
```

## 📚 常用组件列表

推荐先添加以下基础组件：

```bash
npx shadcn-vue@latest add button
npx shadcn-vue@latest add card
npx shadcn-vue@latest add input
npx shadcn-vue@latest add label
npx shadcn-vue@latest add table
npx shadcn-vue@latest add dialog
npx shadcn-vue@latest add select
npx shadcn-vue@latest add dropdown-menu
npx shadcn-vue@latest add badge
npx shadcn-vue@latest add alert
```

## 🔗 相关链接

- [shadcn-vue 官方文档](https://www.shadcn-vue.com/)
- [组件示例](https://www.shadcn-vue.com/docs/components/accordion)
- [Tailwind CSS 文档](https://tailwindcss.com/docs)
- [Radix Vue](https://www.radix-vue.com/) - 底层无样式组件库

## 💡 使用工具函数

我们提供了 `cn()` 工具函数来合并 class 名称：

```vue
<script setup>
import { cn } from '@/lib/utils'

const buttonClass = cn(
  'px-4 py-2',
  'rounded-md',
  isActive && 'bg-primary',
  isDisabled && 'opacity-50 cursor-not-allowed'
)
</script>

<template>
  <button :class="buttonClass">
    Click me
  </button>
</template>
```

## 🚀 下一步

1. 添加你需要的组件
2. 查看现有的页面和组件，使用新的 shadcn-vue 组件替换旧的样式
3. 享受现代化、美观的 UI！
