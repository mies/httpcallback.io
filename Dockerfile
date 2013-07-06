from base

# Add user that will run the service
RUN useradd -m httpcallback

# Make bin path that will hold the binary
RUN mkdir -p /home/httpcallback/bin

# Make sure the user owns his home
RUN chown -R httpcallback:httpcallback /home/httpcallback

# Add the new release to the container
ADD httpcallback.io" /home/httpcallback/bin

# Make httpcallback.io executable
RUN chown u+x /home/httpcallback/bin/httpcallback.io

# Expose the port
EXPOSE 8000

CMD su - httpcallback /home/httpcallback/bin/httpcallback.io
