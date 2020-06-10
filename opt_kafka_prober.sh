#!/bin/bash
# 生成创建broker的指令：cat kafka-cluster.txt|while read line;do zoo_con=$(echo $line|sed 's/9092/2181/');echo "/home/work/local/kafka_online/bin/kafka-topics.sh --zookeeper ${zoo_con} --create --topic kafka-prober-topic1 --partitions 3 --replication-factor 3";done
topic="kafka-prober-topic1"
#topic="kafka-prober-topic3"
username="test"
passwd="test"
interval=30

function kill_one(){
  echo "stoping kafka-probe in ${mode} mode for ${cluster} cluster"
  ps -ef|grep kafka-prober|grep ${mode}|grep ${cluster}|awk '{print $2}'|xargs -i kill  {}
  is_exist=`check_one`
  if [ "${is_exist}X" != "X" ];then
    echo "stop failed"
  else
    echo "stop success"
    exit
  fi
}
function check_one(){
  #echo "checking kafka-probe in ${mode} mode for ${cluster} cluster"
  proc=$(ps -ef|grep kafka-prober|grep ${mode}|grep ${cluster}|grep -v grep)
  if [ "${proc}X" != "X" ];then
    echo -e "\033[32m${cluster} is runing \033[0m"
  else
    echo -e "\033[31m${cluster} is stop \033[0m"
  #  exit
  fi
}
function exec_start(){
    echo "starting kafka-probe in ${mode} mode for ${cluster} cluster"
    is_exist=`check_one`
    if [ "${is_exist}X" != "X" ];then
        echo "kafka-probe in ${mode} mode for ${cluster} cluster is runing,ignore this option"
        #exit
    fi
    nohup ./kafka-prober -brokers "${brokers}" -interval ${interval} -topics "${topic}" -mode ${mode} -username="${username}" -passwd ${passwd} -cluster ${cluster} -log ${log} -config ${conf} >kafka-prober-${mode}-${cluster}.exec.log &
    is_exist=`check_one`
    if [ "${is_exist}X" != "X" ];then
      echo "start success"
    else
      echo "start failed"
      #exit
    fi
}

function help(){
  echo -e "-m 指定启动的模式，可选项：produce、consume\n-c 集群名\n-o 操作类型，可选项: start、stop、check"
  echo "Example:"
  echo "  以生产者模式启动,探测kafka_test_acl集群"
  echo "  sh $0 -m produce -c kafka_test_acl -o start"
  echo "  以生产者模式启动,探测所有已在kafka-cluster.txt中配置过的集群"
  echo "  sh $0 -m produce -o start"
  echo "  指定配置文件启动，默认的配置文件开启sasl，如果未开启sasl则指定配置文件"
  echo "  sh $0 -m produce -c kafka_test_acl -f ./kafka-prober-noauth.yaml -o start"
  echo "  检查所有的produce服务"
  echo "  sh $0 -m produce -o check"
  exit
}
function check_all(){
  all=$(cat kafka-cluster.txt |cut -d'/' -f2)
  run=$(ps -ef|grep produce|grep -o '\-cluster .*'|awk '{print $2}')
  for line in $all;do
    if echo $run|grep -w $line &>/dev/null;then
      echo -e "\033[32m${line} is runing \033[0m"
    else
      echo -e "\033[31m${line} is stop \033[0m"
    fi
  done
}
while getopts "m:c:o:f:" arg;do
    case $arg in
        m)mode=$OPTARG;;
        c)cluster=$OPTARG;;
        o)opt=$OPTARG;;
        f)conf=$OPTARG;;
        ?)echo "unkonw argument";;
    esac
done

if [ "${mode}X" == "X" ]||[ "${opt}X" == "X" ];then
  help
fi
if [ "${cluster}X" == "X" ];then
  echo "Warning,exec operation for all clusters"
  clusters=`cat kafka-cluster.txt|cut -d'/' -f2`
  # brokers_all=`cat kafka-cluster.txt|cut -d'/' -f1`
else
  clusters=${cluster}
fi
if [ "${conf}X" == "X" ];then
  conf="/home/work/kafka_prober/kafka-prober.yaml"
fi
for cluster in ${clusters};do
    log="/home/work/kafka_prober/logs/kafka-prober-${cluster}.log"
    brokers=`egrep "${cluster}$" kafka-cluster.txt|cut -d'/' -f1`
    if [ "${brokers}X" == "X" ];then
      echo "cluster name ${cluster} not exits"
      exit
    fi
    case ${opt} in
      "start")
        exec_start;;
      "stop")
        kill_one;;
      "check")
        check_one;;
      #"checkall")
      #  check_all;;
      *)
        echo "wrong args for -o"
        help;;
    esac

done
