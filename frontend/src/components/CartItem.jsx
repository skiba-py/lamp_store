import { Box, Flex, Image, Text, IconButton } from '@chakra-ui/react';
import { CloseIcon } from '@chakra-ui/icons';
import { Link } from 'react-router-dom';

export default function CartItem({ product, quantity, onRemove }) {
  return (
    <Flex align="center" borderWidth={1} borderRadius="md" p={3} mb={3} bg="white">
      <Link to={`/product/${product.id}`} style={{ display: 'flex', alignItems: 'center' }}>
        <Image src={product.image} alt={product.name} boxSize="60px" objectFit="cover" borderRadius="md" mr={4} />
      </Link>
      <Box flex={1}>
        <Link to={`/product/${product.id}`} style={{ textDecoration: 'none', color: 'inherit' }}>
          <Text fontWeight="bold">{product.name}</Text>
        </Link>
        <Text color="gray.500">{product.price} $</Text>
        <Text color="gray.600">Кол-во: {quantity}</Text>
      </Box>
      <IconButton icon={<CloseIcon />} size="sm" onClick={onRemove} aria-label="Удалить" />
    </Flex>
  );
}
