from flask import Flask, request, jsonify

app = Flask(__name__)

products = []

@app.route("/products", methods=["GET"])
def get_products():
    return jsonify(products)

@app.route("/products", methods=["POST"])
def create_product():
    product = request.json
    products.append(product)
    return jsonify(product)

if __name__ == "__main__":
    app.run(host="0.0.0.0", port=3002)
