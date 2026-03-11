output "instance_id" {
  value = aws_instance.production_server.id
}

output "public_ip" {
  value = aws_instance.production_server.public_ip
}
