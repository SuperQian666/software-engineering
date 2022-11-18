# software-engineering
软件工程实验

本项目基于windows gui运行，使用是直接运行hello.exe即可

构建本项目步骤：
  - 拉取本项目代码
  - 拉取go.mod下所有依赖
  - go build -o hello.exe -ldflags="-H windowsgui"
