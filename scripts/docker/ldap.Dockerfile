FROM osixia/openldap:stable

ENV LDAP_ORGANISATION="wyGOps Dev Environment" \
    LDAP_DOMAIN="wygops.internal"

COPY bootstrap.ldif /container/service/slapd/assets/config/bootstrap/ldif/50-bootstrap.ldif
COPY memberof.ldiff /container/service/slapd/assets/config/bootstrap/ldif/60-memberof.ldiff