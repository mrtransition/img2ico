```markdown
# img2ico

将常见图片（PNG、JPEG、GIF、BMP）一键转换为 Windows ICO 图标文件的命令行工具。

## 功能特性

- 支持输入格式：PNG、JPEG、GIF、BMP
- 自动生成多尺寸 ICO（默认 16×16、32×32、48×48、256×256）
- 自定义输出尺寸（1~256 任意尺寸）
- 支持批量转换与通配符匹配（如 `*.png`）
- 保留透明通道
- 跨平台支持（Windows / macOS / Linux）

## 安装

### 方式一：从源码构建

```bash
git clone https://github.com/mrtransition/img2ico.git
cd img2ico
go build -o img2ico
```

### 方式二：直接下载二进制

前往 [Releases](https://github.com/mrtransition/img2ico/releases) 页面下载对应平台的二进制文件。

## 使用方法

### 基本转换

```bash
# 将 logo.png 转为 logo.ico（默认尺寸：16,32,48,256）
img2ico logo.png

# 指定输出文件路径
img2ico logo.png -o favicon.ico

# 指定输出目录
img2ico logo.png -o ./icons/
```

### 自定义尺寸

```bash
# 只生成 16×16 和 32×32
img2ico logo.png -s 16,32

# 生成所有标准尺寸（16,24,32,48,64,72,96,128,256）
img2ico logo.png --all
```

### 批量转换

```bash
# 转换当前目录所有 PNG 文件
img2ico "*.png"

# 转换多个文件，输出到指定目录
img2ico pic1.png pic2.jpg pic3.gif -o ./icons/

# 递归批量处理（配合通配符）
img2ico "./assets/**/*.png" -o ./output/ --overwrite
```

### 其他选项

```bash
# 覆盖已存在的 ICO 文件
img2ico logo.png --overwrite

# 显示详细转换信息
img2ico logo.png --verbose
```

## 命令行参数

| 参数 | 简写 | 说明 | 默认值 |
|------|------|------|--------|
| `--output` | `-o` | 输出文件（单文件模式）或输出目录（多文件模式） | 当前目录 |
| `--sizes` | `-s` | 尺寸列表，逗号分隔 | `16,32,48,256` |
| `--all` | `-m` | 生成所有标准尺寸（覆盖 `--sizes`） | `false` |
| `--overwrite` | `-w` | 覆盖已有文件 | `false` |
| `--verbose` | `-v` | 显示详细日志 | `false` |

## 示例输出

```bash
$ img2ico app.png -v
Input files: [app.png]
Sizes: [16 32 48 256]
Converting app.png (256x256) to app.ico with sizes [16 32 48 256]
  - size 16 encoded (958 bytes)
  - size 32 encoded (2532 bytes)
  - size 48 encoded (4433 bytes)
  - size 256 encoded (19744 bytes)
Created: app.ico (contains 4 images)
All conversions completed successfully.
```

## 技术细节

- 内部使用标准库 `image` 解码，支持 PNG、JPEG、GIF、BMP。
- 缩放算法采用 Lanczos 滤波（`github.com/disintegration/imaging`），保证高质量。
- ICO 编码为**纯手工实现**，符合 Windows ICO 格式规范，直接写入 PNG 数据（Windows Vista 以上原生支持）。
- 多尺寸图片会全部写入同一个 ICO 文件，完全兼容 Windows 资源管理器及各类应用程序。

## 常见问题

**Q: 为什么在资源管理器中预览只看到一个尺寸？**  
A: 这是正常现象，Windows 缩略图默认只显示第一个尺寸（通常是最小的）。您可以在文件**属性 → 详细信息**中查看所有包含的尺寸，或使用 [Icon Viewer](https://www.nirsoft.net/utils/iconsext.html) 等工具验证。

**Q: 可以生成 256×256 以上的图标吗？**  
A: ICO 格式标准最大尺寸为 256×256。更大的尺寸不会被系统识别，本工具会拒绝超出范围的值。

**Q: 支持 CUR（光标文件）格式吗？**  
A: 目前仅支持 ICO。如果需要光标功能，欢迎提交 Issue。

## 许可证

MIT License

## 贡献

欢迎提交 Issue 和 Pull Request。

---
Made with ❤️ by [deepseek]
```
