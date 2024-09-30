FROM rust:1.81.0-bookworm AS build

ARG APP_NAME=kbservice
ARG appVersion
ARG commitHash
ARG buildDate

ENV VERSION=$appVersion
ENV COMMIT_HASH=$commitHash
ENV BUILD_DATE=$buildDate

WORKDIR /build

COPY Cargo.lock Cargo.toml ./
RUN mkdir src \
    && echo "// dummy file" > src/lib.rs \
    && cargo build --release

COPY src src
RUN touch src/lib.rs
RUN cargo build --locked --release
RUN cp ./target/release/$APP_NAME /bin/server

FROM debian:bookworm-slim AS final
COPY --from=build /bin/server /bin/
EXPOSE 3030
CMD ["/bin/server"]