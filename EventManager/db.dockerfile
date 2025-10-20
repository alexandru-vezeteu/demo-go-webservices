FROM mysql:8.0.43-debian

COPY ./db /scripts

# ceva de genul trb sa schimb probabil sa log in and stuff
RUN "/bin/bash" /scripts/create.sql
RUN "/bin/bash" /scripts/populate.sql