
NAME = flutter-image

launch_apk: build apk
# to launch target web-server
launch_web:	build web 

build:
	@echo "Building Docker image ..."	
	@docker build -t $(NAME) .
	@echo "Docker image built successfully."
	
run:
	@echo "Starting container for Web Dev..."
# Added '-' to ignore exit error and fixed '--rm'
	-@docker run --rm -it \
		-p 8080:8080 \
		-v $(PWD)/mobilePiscine:/app \
		$(NAME) bash

# Option 1: Build Android apk
apk:
	@echo "Flutter building Android APK..."
	-@docker run --rm -it \
	  -v $(PWD)/mobilePiscine:/app \
	  -w /app \
	  $(NAME) bash -lc '\
	  	cd mobileModule00/ex00 && \
	  	flutter clean && \
	  	flutter pub get && \
	  	flutter build apk --release'
	@echo "Build finished..."
	@echo "Check mobilePiscine/mobileModule00/ex00/build/app/outputs/flutter-apk/"

# Option 2: After you run make launch_Web, go to browser to load app, http://localhost:8080
web:
	@echo "Flutter run app on web.."
# Added '-' to ignore exit error even if flutter crashes
	-@docker run --rm -it \
	  -p 8080:8080 \
	  -p 5555:5555 \
	  -v $(PWD)/mobilePiscine:/app \
	  -w /app \
	  $(NAME) bash -lc '\
	    flutter config --enable-web && \
		cd mobileModule00/ex00 && \
	    flutter pub get && \
		mkdir -p /home/developer/build_cache && \
		rm -rf build && \
		ln -s /home/developer/build_cache build && \
	    flutter run -d web-server \
	      --web-port 8080 \
	      --web-hostname 0.0.0.0 \
	      --host-vmservice-port 5555'
	@echo "Visit http://localhost:8080 to load app."
	
clean:
	@echo "Removing built images..."
	@docker image rm -f $(NAME) >/dev/null 2>&1 || true
	@echo "Image removed."

# container user and host users are different. To have permissions to delete files from host.	
chmod:
	@ID=$$(docker ps -q --filter ancestor=$(NAME)); \
	if [ -n "$$ID" ]; then \
		docker exec -u 0 $$ID sh -c "chmod -R a+rwx /app/mobileModule00"; \
		echo "Permissions fixed for container $$ID"; \
	else \
		echo "No running container found."; \
	fi
	
prune:	clean
	@echo "Cleaning up Docker build cache..."
	@docker system prune -af || true
	@docker volume prune -f
	@echo "System and volume pruning done."

.PHONY: launch_apk launch_web build run web chmod clean prune

