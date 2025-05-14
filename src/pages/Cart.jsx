import { Box, Heading, Text } from '@chakra-ui/react';
import CartList from '../components/CartList';
import { products } from '../data';

// Моковые данные корзины
const cartItems = [
  { product: products[0], quantity: 2 },
  { product: products[1], quantity: 1 },
];

export default function Cart() {
  const handleRemove = (productId) => {
    alert(`Удалить товар с id: ${productId}`);
  };

  const total = cartItems.reduce((sum, item) => sum + item.product.price * item.quantity, 0);

  return (
    <Box id="cart" w="100%" minH="100vh" maxW="600px">
      <Heading as="h1" size="lg" mb={6}>Корзина</Heading>
      <CartList items={cartItems} onRemove={handleRemove} />
      <Box mt={6} textAlign="right">
        <Text fontWeight="bold" fontSize="xl">Итого: {total} $</Text>
      </Box>
    </Box>
  );
}
