# Go Asistant

### Porcupine hotword detection (deprecated)

由於不支援 Raspberry pi 產生辨識 model，所以改採 SnowBoy 開發。

#### 使用

安裝依賴
```
cp resources/porcupine/lib/linux/x86_64/* /usr/local/lib64/
cp resources/porcupine/include/* /usr/local/include

export LD_LIBRARY_PATH=/usr/local/lib64
```

執行
```
sudo chmod +x ./start.sh
./start.sh
```
