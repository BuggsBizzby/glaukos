FROM python:3 AS mitmproxy

USER root

WORKDIR /root
RUN pip install mitmproxy
RUN timeout 1 mitmdump || true

RUN mkdir -p /data && chmod 777 /data

CMD ["mitmdump", "-w", "+/data/web_traffic.dmp"]


FROM kasmweb/chromium:1.14.0-rolling AS chromium

USER root

RUN locale-gen en_US.UTF-8
RUN update-locale

RUN apt-get update
RUN apt-get install -y libnss3-tools

WORKDIR /root

COPY --from=mitmproxy /root/.mitmproxy/mitmproxy-ca-cert.pem /usr/local/share/ca-certificates/mitmproxy-ca-cert.crt
RUN update-ca-certificates

WORKDIR /usr/share/kasmvnc/www/app/images/icons
RUN rm *
COPY favicon.png .

WORKDIR /usr/share/kasmvnc/www
COPY vnc_visual_fixes.py .
RUN python3 /usr/share/kasmvnc/www/vnc_visual_fixes.py

WORKDIR /usr/share/kasmvnc/www/dist/images
RUN rm *

USER kasm-user

RUN set -x

RUN echo 'rm -r /home/kasm-user/.pki/nssdb || true' >> /home/kasm-default-profile/.bashrc
RUN echo 'mkdir -p /home/kasm-user/.pki/nssdb' >> /home/kasm-default-profile/.bashrc
RUN echo 'certutil -N -d sql:/home/kasm-user/.pki/nssdb --empty-password < /dev/null' >> /home/kasm-default-profile/.bashrc
RUN echo 'certutil -d sql:/home/kasm-user/.pki/nssdb -A -t "C,," -n "mitmproxy2" -i /usr/local/share/ca-certificates/mitmproxy-ca-cert.crt' >> /home/kasm-default-profile/.bashrc
