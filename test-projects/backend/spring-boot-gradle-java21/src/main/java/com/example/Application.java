package com.example;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.web.bind.annotation.*;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;

import java.time.Instant;
import java.util.*;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.concurrent.atomic.AtomicLong;

@SpringBootApplication
@RestController
public class Application {

    private final List<Item> items = Collections.synchronizedList(new ArrayList<>());
    private final AtomicInteger nextId = new AtomicInteger(3);
    private final AtomicLong requestCount = new AtomicLong(0);
    private final long startTime = System.currentTimeMillis();

    public Application() {
        items.add(new Item(1, "Sample Item 1", "First sample item", Instant.now().toString()));
        items.add(new Item(2, "Sample Item 2", "Second sample item", Instant.now().toString()));
    }

    public static void main(String[] args) {
        SpringApplication.run(Application.class, args);
    }

    @GetMapping("/")
    public Map<String, Object> getRoot() {
        requestCount.incrementAndGet();
        Map<String, Object> response = new HashMap<>();
        response.put("message", "Spring Boot Test API");
        response.put("version", "1.0.0");
        response.put("requestCount", requestCount.get());
        response.put("timestamp", Instant.now().toString());
        response.put("endpoints", Arrays.asList(
            "GET /",
            "GET /health",
            "GET /items",
            "POST /items",
            "GET /items/{id}",
            "DELETE /items/{id}"
        ));
        return response;
    }

    @GetMapping("/health")
    public Map<String, Object> getHealth() {
        requestCount.incrementAndGet();
        long uptime = (System.currentTimeMillis() - startTime) / 1000;
        Map<String, Object> response = new HashMap<>();
        response.put("status", "healthy");
        response.put("uptime", uptime + " seconds");
        response.put("timestamp", Instant.now().toString());
        return response;
    }

    @GetMapping("/items")
    public Map<String, Object> getItems() {
        requestCount.incrementAndGet();
        Map<String, Object> response = new HashMap<>();
        response.put("success", true);
        response.put("count", items.size());
        response.put("data", items);
        return response;
    }

    @PostMapping("/items")
    public ResponseEntity<Map<String, Object>> createItem(@RequestBody CreateItemRequest request) {
        requestCount.incrementAndGet();

        if (request.getName() == null || request.getName().trim().isEmpty()) {
            Map<String, Object> error = new HashMap<>();
            error.put("success", false);
            error.put("error", "Name is required");
            return ResponseEntity.badRequest().body(error);
        }

        Item newItem = new Item(
            nextId.getAndIncrement(),
            request.getName(),
            request.getDescription() != null ? request.getDescription() : "",
            Instant.now().toString()
        );
        items.add(newItem);

        Map<String, Object> response = new HashMap<>();
        response.put("success", true);
        response.put("data", newItem);
        return ResponseEntity.status(HttpStatus.CREATED).body(response);
    }

    @GetMapping("/items/{id}")
    public ResponseEntity<Map<String, Object>> getItem(@PathVariable int id) {
        requestCount.incrementAndGet();
        Optional<Item> item = items.stream().filter(i -> i.getId() == id).findFirst();

        if (item.isEmpty()) {
            Map<String, Object> error = new HashMap<>();
            error.put("success", false);
            error.put("error", "Item not found");
            return ResponseEntity.status(HttpStatus.NOT_FOUND).body(error);
        }

        Map<String, Object> response = new HashMap<>();
        response.put("success", true);
        response.put("data", item.get());
        return ResponseEntity.ok(response);
    }

    @DeleteMapping("/items/{id}")
    public ResponseEntity<Map<String, Object>> deleteItem(@PathVariable int id) {
        requestCount.incrementAndGet();
        Optional<Item> item = items.stream().filter(i -> i.getId() == id).findFirst();

        if (item.isEmpty()) {
            Map<String, Object> error = new HashMap<>();
            error.put("success", false);
            error.put("error", "Item not found");
            return ResponseEntity.status(HttpStatus.NOT_FOUND).body(error);
        }

        items.remove(item.get());
        Map<String, Object> response = new HashMap<>();
        response.put("success", true);
        response.put("message", "Item deleted successfully");
        response.put("data", item.get());
        return ResponseEntity.ok(response);
    }

    // Inner classes
    static class Item {
        private int id;
        private String name;
        private String description;
        private String createdAt;

        public Item(int id, String name, String description, String createdAt) {
            this.id = id;
            this.name = name;
            this.description = description;
            this.createdAt = createdAt;
        }

        public int getId() { return id; }
        public String getName() { return name; }
        public String getDescription() { return description; }
        public String getCreatedAt() { return createdAt; }
    }

    static class CreateItemRequest {
        private String name;
        private String description;

        public String getName() { return name; }
        public void setName(String name) { this.name = name; }
        public String getDescription() { return description; }
        public void setDescription(String description) { this.description = description; }
    }
}
