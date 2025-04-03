----

## 使用 `tmux` 创建多窗口会话

`tmux` 是类似 `screen` 的工具，功能更丰富，推荐用来运行持久化服务。  

1. 启动 `tmux` 会话：  

```bash
tmux new -s emqx-pg-bridge
```

2. 在 `tmux` 窗口中运行程序：  

```bash
./emqx-pg-bridge -config ./config.yaml

/opt/myapp/apps/emqx-pg-bridge/emqx-pg-bridge -config /opt/myapp/apps/emqx-pg-bridge/config.yaml
```

3. **退出但保持程序运行：**  
   按 `Ctrl + B`，然后按 `D`（表示 detach）。  

4. **重新连接到会话：**  

```bash
tmux attach -t emqx-pg-bridge
```

5. **关闭会话（终止程序）：**  
   在会话中按 `Ctrl + C` 终止程序，或者在外部直接关闭会话：  

```bash
tmux kill-session -t emqx-pg-bridge
```

---