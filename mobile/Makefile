
NAME = flutter-image

launch: build run
app:	build web 

build:
	@echo "Building Docker image ..."	
	@docker build -t $(NAME) .
	@echo "Docker image built successfully."
	
run:
	@echo "Starting container for Web Dev..."
	@docker run -it --rm \
		-p 8080:8080 \
		-v $(PWD)/mobilePiscine:/app \
		$(NAME) bash
		
# After you run make app, go to browser to load app, http://localhost:8080
web:
	@echo "Flutter run app on web.."
	@docker run --rm -it \
	  -p 8080:8080 \
	  -p 5555:5555 \
	  -v $(PWD)/mobilePiscine:/app \
	  -w /app \
	  $(NAME) bash -lc '\
	    flutter config --enable-web && \
	    cd mobileModule00/ex00 && \
	    flutter run -d web-server \
	      --web-port 8080 \
	      --web-hostname 0.0.0.0 \
	      --host-vmservice-port 5555'
	@echo "Visit http://localhost:8080 to load app."
	
clean:
	@echo "Removing built images..."
	@docker image rm -f $(NAME) || true
	@echo "Image removed."

# container user and host users are different. To have permissions to delete files from host.	
chmod:
	@docker exec -u 0 6db4fac85471 sh -c "chmod -R a+rwx /app/mobileModule00"
	@echo "chmod -R a+rwx /app/mobileModule00 done."
	
prune:	clean
	@echo "Cleaning up Docker build cache..."
	@docker system prune -af || true
	@docker volume prune -f
	@echo "System and volume pruning done."

.PHONY: launch app build run web chmod clean prune

