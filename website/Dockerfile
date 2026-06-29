FROM hub.intra.mlamp.cn/library/nginx:1.17 AS runtime
COPY dist/ /usr/share/nginx/html/
RUN chmod -R 755 /usr/share/nginx/html && \
    printf 'server {\n  listen 80;\n  root /usr/share/nginx/html;\n  index index.html;\n  location / {\n    try_files $uri $uri/ /index.html;\n  }\n}\n' > /etc/nginx/conf.d/default.conf
EXPOSE 80
