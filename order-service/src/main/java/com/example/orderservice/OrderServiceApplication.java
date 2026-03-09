package com.example.orderservice;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.web.bind.annotation.*;

import java.util.ArrayList;
import java.util.List;

@SpringBootApplication
@RestController
public class OrderServiceApplication {

    List<String> orders = new ArrayList<>();

    public static void main(String[] args) {
        SpringApplication.run(OrderServiceApplication.class, args);
    }

    @GetMapping("/orders")
    public List<String> getOrders() {
        return orders;
    }

    @PostMapping("/orders")
    public String createOrder(@RequestBody String order) {
        orders.add(order);
        return order;
    }
}
