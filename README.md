# klog

a simple view k8s logging tool, it depend on kubectl tool.

<br>

# usage

> go get github.com/zhufuyi/klog

go into the folder, run command **go build**, and move binary to your envionment.

<br>

view the first pod log:

> klog --name=account

view spec pod log:

> klog --name=account --podNO=3
