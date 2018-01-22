# Mac下配置VirtualBox共享文件夹

## 安装VirtualBox和扩展插件
1. 下载最新版安装包和插件
https://www.virtualbox.org/wiki/Downloads
2. 对应版本扩展文件:Oracle_VM_VirtualBox_Extension_Pack-[version]-119785.vbox-extpack
3. 在VirtualBox -> 偏好设置 -> 扩展 点旁边“+”号，选中上述文件安装
PS：只有对应版本可以安装成功

## 在VirtualBox中安装Ubuntu
操作系统版本要求：16.04 + 

## 虚拟机配置
1. 在虚拟机启动窗口顶部目录中； 设备 -> 安装增强功能，此时，VirtualBox虚拟机控制界面中，存储 -> 控制器IDE中出现VBoxGuestAdditions.iso
2. 在虚拟机控制界面中，共享文件夹中配置共享文件夹目录
3. 在虚拟机操作系统中执行挂载命令：
mount -t vboxsf dirname localpath
