package com.example

import org.springframework.boot.autoconfigure.SpringBootApplication
import org.springframework.boot.runApplication
import org.springframework.web.bind.annotation.*
import org.springframework.http.HttpStatus
import org.springframework.http.ResponseEntity
import java.time.Instant
import java.util.concurrent.atomic.AtomicInteger
import java.util.concurrent.atomic.AtomicLong

data class Item(
    val id: Int,
    val name: String,
    val description: String,
    val createdAt: String
)

data class CreateItemRequest(
    val name: String?,
    val description: String?
)

@SpringBootApplication
@RestController
class Application {

    private val items = mutableListOf(
        Item(1, "Sample Item 1", "First sample item", Instant.now().toString()),
        Item(2, "Sample Item 2", "Second sample item", Instant.now().toString())
    )
    private val nextId = AtomicInteger(3)
    private val requestCount = AtomicLong(0)
    private val startTime = System.currentTimeMillis()

    @GetMapping("/")
    fun getRoot(): Map<String, Any> {
        requestCount.incrementAndGet()
        return mapOf(
            "message" to "Kotlin Spring Boot Test API",
            "version" to "1.0.0",
            "requestCount" to requestCount.get(),
            "timestamp" to Instant.now().toString(),
            "endpoints" to listOf(
                "GET /",
                "GET /health",
                "GET /items",
                "POST /items",
                "GET /items/{id}",
                "DELETE /items/{id}"
            )
        )
    }

    @GetMapping("/health")
    fun getHealth(): Map<String, Any> {
        requestCount.incrementAndGet()
        val uptime = (System.currentTimeMillis() - startTime) / 1000
        return mapOf(
            "status" to "healthy",
            "uptime" to "$uptime seconds",
            "timestamp" to Instant.now().toString()
        )
    }

    @GetMapping("/items")
    fun getItems(): Map<String, Any> {
        requestCount.incrementAndGet()
        return mapOf(
            "success" to true,
            "count" to items.size,
            "data" to items
        )
    }

    @PostMapping("/items")
    fun createItem(@RequestBody request: CreateItemRequest): ResponseEntity<Map<String, Any>> {
        requestCount.incrementAndGet()

        if (request.name.isNullOrBlank()) {
            return ResponseEntity.badRequest().body(
                mapOf(
                    "success" to false,
                    "error" to "Name is required"
                )
            )
        }

        val newItem = Item(
            id = nextId.getAndIncrement(),
            name = request.name,
            description = request.description ?: "",
            createdAt = Instant.now().toString()
        )
        items.add(newItem)

        return ResponseEntity.status(HttpStatus.CREATED).body(
            mapOf(
                "success" to true,
                "data" to newItem
            )
        )
    }

    @GetMapping("/items/{id}")
    fun getItem(@PathVariable id: Int): ResponseEntity<Map<String, Any>> {
        requestCount.incrementAndGet()
        val item = items.find { it.id == id }

        return if (item != null) {
            ResponseEntity.ok(
                mapOf(
                    "success" to true,
                    "data" to item
                )
            )
        } else {
            ResponseEntity.status(HttpStatus.NOT_FOUND).body(
                mapOf(
                    "success" to false,
                    "error" to "Item not found"
                )
            )
        }
    }

    @DeleteMapping("/items/{id}")
    fun deleteItem(@PathVariable id: Int): ResponseEntity<Map<String, Any>> {
        requestCount.incrementAndGet()
        val item = items.find { it.id == id }

        return if (item != null) {
            items.remove(item)
            ResponseEntity.ok(
                mapOf(
                    "success" to true,
                    "message" to "Item deleted successfully",
                    "data" to item
                )
            )
        } else {
            ResponseEntity.status(HttpStatus.NOT_FOUND).body(
                mapOf(
                    "success" to false,
                    "error" to "Item not found"
                )
            )
        }
    }
}

fun main(args: Array<String>) {
    runApplication<Application>(*args)
}
