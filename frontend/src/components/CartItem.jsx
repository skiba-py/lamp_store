import { Box, Image, Text, IconButton, Button, Flex } from '@chakra-ui/react';
import { FaTrash } from 'react-icons/fa';
import { Link } from 'react-router-dom';

export default function CartItem({ item, onRemove, onUpdateQuantity, updating }) {
  return (
    <Flex align="center" borderWidth={1} borderRadius="md" p={3} mb={3} bg="white">
      <Box w="80px" h="80px" mr={4}>
        <Image
          src={item.product && item.product.image ? item.product.image : undefined}
          alt={item.product ? item.product.name : ''}
          objectFit="cover"
          w="100%"
          h="100%"
          borderRadius="md"
        />
      </Box>
      <Box flex={1}>
        <Link to={`/product/${item.product ? item.product.id : ''}`} style={{ textDecoration: 'none', color: 'inherit' }}>
          <Text fontWeight="bold">{item.product ? item.product.name : 'Товар не найден'}</Text>
        </Link>
        <Text color="gray.500">{item.product && item.product.price ? item.product.price : 0} ₽</Text>
        <Flex align="center" gap={2} mt={2}>
          <Button size="sm" onClick={() => onUpdateQuantity(item.id, -1)} disabled={updating || item.quantity <= 1}>-</Button>
          <Text mx={2}>{item.quantity}</Text>
          <Button size="sm" onClick={() => onUpdateQuantity(item.id, 1)} disabled={updating}>+</Button>
        </Flex>
      </Box>
      <IconButton
        colorScheme="red"
        icon={<FaTrash />}
        onClick={() => onRemove(item.id)}
        size="sm"
        ml={2}
      >Удалить</IconButton>
    </Flex>
  );
}
