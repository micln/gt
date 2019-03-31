# gt

缘起：内部的 Go base 包每天都会更新很多次，如果不打 tag，会出来很多不兼容的 master-0.0.0-xxx。

Go Module 建议每次更新都使用语义化版本，这个工具就是来帮你简化这个工作。 

## Install

    go get github.com/micln/gt

## Usage

```shell
# 打一个 fix/patch tag
gt -fix

# 打一个 feature tag
gt -ft

# help
gt -h
```
