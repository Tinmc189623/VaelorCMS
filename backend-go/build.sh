#!/bin/bash

# VaelorCMS - 多平台编译脚本 (Bash)
# 支持 Windows、Linux、macOS 编译

set -e

# 项目配置
PROJECT_NAME="vaelorcms"
MAIN_PACKAGE="cmd/server"
OUTPUT_DIR="bin"
CLEAN=false
ALL=false

# 目标平台配置
TARGETS=(
    "windows amd64 .exe"
    "windows 386 .exe"
    "linux amd64"
    "linux 386"
    "linux arm64"
    "linux arm"
    "darwin amd64"
    "darwin arm64"
)

# 解析命令行参数
while [[ "$#" -gt 0 ]]; do
    case $1 in
        -o|--output) OUTPUT_DIR="$2"; shift ;;
        -c|--clean) CLEAN=true ;;
        -a|--all) ALL=true ;;
        -h|--help) 
            echo "VaelorCMS - 多平台编译脚本"
            echo ""
            echo "用法: $0 [选项]"
            echo ""
            echo "选项:"
            echo "  -o, --output <目录>    输出目录 (默认: bin)"
            echo "  -c, --clean           清理输出目录"
            echo "  -a, --all             编译所有平台"
            echo "  -h, --help            显示帮助信息"
            echo ""
            exit 0
            ;;
        *) echo "未知选项: $1"; exit 1 ;;
    esac
    shift
done

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# 清理输出目录
if [ "$CLEAN" = true ] && [ -d "$OUTPUT_DIR" ]; then
    echo -e "${YELLOW}正在清理输出目录: $OUTPUT_DIR${NC}"
    rm -rf "$OUTPUT_DIR"
fi

# 创建输出目录
mkdir -p "$OUTPUT_DIR"

# 获取当前平台
CURRENT_OS=$(uname -s | tr '[:upper:]' '[:lower:]')
if [ "$CURRENT_OS" = "darwin" ]; then
    CURRENT_OS="darwin"
elif [[ "$CURRENT_OS" == *"mingw"* ]] || [[ "$CURRENT_OS" == *"msys"* ]] || [[ "$CURRENT_OS" == *"cygwin"* ]]; then
    CURRENT_OS="windows"
else
    CURRENT_OS="linux"
fi

CURRENT_ARCH=$(uname -m)
case "$CURRENT_ARCH" in
    x86_64) CURRENT_ARCH="amd64" ;;
    i686|i386) CURRENT_ARCH="386" ;;
    aarch64|arm64) CURRENT_ARCH="arm64" ;;
    armv*) CURRENT_ARCH="arm" ;;
esac

echo -e "${CYAN}========================================${NC}"
echo -e "${CYAN}  VaelorCMS - 多平台编译脚本${NC}"
echo -e "${CYAN}========================================${NC}"
echo ""

# 如果没有指定 --all，只编译当前平台
if [ "$ALL" = false ]; then
    TARGETS=("$CURRENT_OS $CURRENT_ARCH")
    echo -e "${GREEN}仅编译当前平台: $CURRENT_OS/$CURRENT_ARCH${NC}"
else
    echo -e "${GREEN}编译所有平台...${NC}"
fi
echo ""

SUCCESS_COUNT=0
FAIL_COUNT=0

for TARGET in "${TARGETS[@]}"; do
    read -r OS ARCH EXT <<< "$TARGET"
    OUTPUT_FILE="$OUTPUT_DIR/${PROJECT_NAME}-${OS}-${ARCH}${EXT}"
    
    echo -e "${CYAN}正在编译: $OS/$ARCH${NC}"
    
    START=$(date +%s.%N)
    
    if GOOS="$OS" GOARCH="$ARCH" CGO_ENABLED=0 go build -ldflags "-s -w" -o "$OUTPUT_FILE" "$MAIN_PACKAGE"; then
        END=$(date +%s.%N)
        DURATION=$(echo "$END - $START" | bc)
        DURATION=$(printf "%.2f" $DURATION)
        
        if [ -f "$OUTPUT_FILE" ]; then
            FILE_SIZE=$(du -k "$OUTPUT_FILE" | cut -f1)
            echo -e "${GREEN}  ✓ 成功! ${FILE_SIZE} KB, 耗时 ${DURATION}s${NC}"
            SUCCESS_COUNT=$((SUCCESS_COUNT + 1))
        else
            echo -e "${RED}  ✗ 失败: 输出文件未生成${NC}"
            FAIL_COUNT=$((FAIL_COUNT + 1))
        fi
    else
        echo -e "${RED}  ✗ 失败: 编译错误${NC}"
        FAIL_COUNT=$((FAIL_COUNT + 1))
    fi
done

echo ""
echo -e "${CYAN}========================================${NC}"
echo -e "${CYAN}  编译完成!${NC}"
echo -e "${CYAN}========================================${NC}"
echo -e "${GREEN}  成功: $SUCCESS_COUNT${NC}"
if [ "$FAIL_COUNT" -gt 0 ]; then
    echo -e "${RED}  失败: $FAIL_COUNT${NC}"
else
    echo -e "${NC}  失败: $FAIL_COUNT${NC}"
fi
echo -e "${CYAN}  输出目录: $(cd "$OUTPUT_DIR" && pwd)${NC}"
echo -e "${CYAN}========================================${NC}"

if [ "$FAIL_COUNT" -gt 0 ]; then
    exit 1
fi
