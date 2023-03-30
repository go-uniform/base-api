FROM scratch
ADD build /service
ENTRYPOINT ["/service"]
HEALTHCHECK CMD ["/service command:ping"]
CMD ["run"]