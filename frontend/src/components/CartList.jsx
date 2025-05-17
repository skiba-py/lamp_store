import CartItem from './CartItem';
import { Box, Text } from '@chakra-ui/react';

export default function CartList({ items, onRemove }) {
  if (items.length === 0) {
    return <Text color="gray.500">Корзина пуста</Text>;
  }
  return (
    <Box>
      {items.map(item => (
        <CartItem key={item.product.id} product={item.product} quantity={item.quantity} onRemove={() => onRemove(item.product.id)} />
      ))}
    </Box>
  );
}
