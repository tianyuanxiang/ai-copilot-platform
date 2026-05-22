FROM nginx:alpine

# 直接复制本地构建好的前端产物
COPY dist /usr/share/nginx/html

# nginx 配置通过 ConfigMap 挂载
RUN rm -f /etc/nginx/conf.d/default.conf

EXPOSE 8356

CMD ["nginx", "-g", "daemon off;"]
