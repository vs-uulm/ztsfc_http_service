FROM ubuntu:latest                                                                      

RUN touch /conf.yml                                                                     
RUN mkdir /basic_auth                                                                   
RUN mkdir /certs                                                                        
RUN mkdir -p /etc/letsencrypt/live/                                                     
RUN mkdir -p /etc/letsencrypt/archive                                                   

EXPOSE 443/tcp                                                                          
# EXPOSE 636/tcp                                                                          
# EXPOSE 389/tcp                                                                          

ADD main /main                                                                          

CMD /main -m pep