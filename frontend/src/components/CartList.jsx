import CartItem from './CartItem';
import { Box, Text } from '@chakra-ui/react';

export default function CartList({ items, onRemove, onUpdateQuantity, updating }) {
  if (!items || items.length === 0) {
    return <Text color="gray.500">Корзина пуста</Text>;
  }
  return (
    <Box>
      {items
        .filter(item => item && item.product)
        .map(item => (
          <CartItem
            key={item.id}
            item={item}
            onRemove={onRemove}
            onUpdateQuantity={onUpdateQuantity}
            updating={updating}
          />
        ))}
    </Box>
  );
}
