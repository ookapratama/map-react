import React, { useState, useEffect, useRef } from "react";
import { Stage, Layer, Rect, Text, Group } from "react-konva";

const ParkingMap = ({ spots, bookings, onSpotClick }) => {
  const containerRef = useRef(null);
  const [dimensions, setDimensions] = useState({
    width: 800,
    height: 400,
  });

  useEffect(() => {
    const handleResize = () => {
      if (containerRef.current) {
        setDimensions({
          width: Math.min(800, containerRef.current.offsetWidth - 32),
          height: 400,
        });
      }
    };

    handleResize();
    window.addEventListener("resize", handleResize);
    return () => window.removeEventListener("resize", handleResize);
  }, []);

  const getSpotStatus = (spotId) => {
    const activeBooking = bookings.find(
      (b) => b.spotId === spotId && b.status === "active",
    );
    return activeBooking ? "occupied" : "available";
  };

  const scale = dimensions.width / 800;

  return (
    <div ref={containerRef} className="w-full flex justify-center">
      <Stage
        width={dimensions.width}
        height={400 * scale}
        scaleX={scale}
        scaleY={scale}
      >
        <Layer>
          {/* Parking Floor */}
          <Rect
            x={0}
            y={0}
            width={800}
            height={400}
            fill="#f8fafc"
            cornerRadius={20}
          />

          {/* Spots */}
          {spots.map((spot) => {
            const status = getSpotStatus(spot.id);
            const isOccupied = status === "occupied";

            return (
              <Group
                key={spot.id}
                onClick={() => onSpotClick(spot)}
                onTap={() => onSpotClick(spot)}
                cursor="pointer"
              >
                {/* Spot Border/Base */}
                <Rect
                  x={spot.x}
                  y={spot.y}
                  width={spot.width}
                  height={spot.height}
                  fill={isOccupied ? "#fee2e2" : "#f0fdf4"}
                  stroke={isOccupied ? "#ef4444" : "#22c55e"}
                  strokeWidth={2}
                  cornerRadius={8}
                />

                {/* Road Marking */}
                <Rect
                  x={spot.x + 10}
                  y={isOccupied ? spot.y + 20 : spot.y + spot.height - 10}
                  width={spot.width - 20}
                  height={isOccupied ? spot.height - 40 : 2}
                  fill={isOccupied ? "#475569" : "#e2e8f0"}
                  cornerRadius={isOccupied ? 12 : 0}
                />

                {/* Spot Label */}
                <Text
                  x={spot.x}
                  y={spot.y + spot.height + 8}
                  width={spot.width}
                  text={spot.label}
                  fontSize={12}
                  fontFamily="Inter"
                  fontStyle="bold"
                  align="center"
                  fill="#94a3b8"
                />
              </Group>
            );
          })}
        </Layer>
      </Stage>
    </div>
  );
};

export default ParkingMap;
