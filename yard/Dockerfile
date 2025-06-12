from node:19-alpine3.16

ARG PRODUCT_VERSION=dev
ENV PRODUCT_VERSION=${PRODUCT_VERSION}
ENV NEXT_PUBLIC_PRODUCT_VERSION=${PRODUCT_VERSION}

RUN mkdir -p /opt/yard
COPY next.config.js package-lock.json package.json postcss.config.js tailwind.config.js start.sh .env.template /opt/yard/

RUN mkdir -p /opt/yard/pages
ADD pages /opt/yard/pages

RUN mkdir -p /opt/yard/components
ADD components /opt/yard/components

RUN mkdir -p /opt/yard/styles
ADD styles /opt/yard/styles

RUN mkdir -p /opt/yard/public
ADD public /opt/yard/public

RUN mkdir -p /opt/yard/utils
ADD utils /opt/yard/utils

WORKDIR /opt/yard
RUN npm install && npm run build

ENTRYPOINT ["npm", "run", "start"]
