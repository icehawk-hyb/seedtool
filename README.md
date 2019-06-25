# 通过seed生成公私钥对,分三种情况处理：
## 1,seed是绝对按照规则过滤多余的空格，并checksum校验通过，使用新的方式生成公私钥
## 2,seed是绝对按照规则过滤多余的空格，并checksum校验通过，使用旧的方式生成公私钥
## 3,seed直接使用客户传入的即可，不需要过滤多余的空格，只能使用旧的方式生成公私钥

# 入参说明：
## 
## Usage of src.exe:
##   -a string
##         account address
##   -c int
##         count (default 100) 生成私钥时index索引取值区间[0-count]
##   -s string
##         seed