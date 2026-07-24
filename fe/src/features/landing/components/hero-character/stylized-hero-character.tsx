'use client'

import { useMemo, useRef } from 'react'
import { useFrame, useThree } from '@react-three/fiber'
import * as THREE from 'three'

import { buildHeroCharacterGroup } from './build-hero-character'

type StylizedHeroCharacterProps = {
  lookStrength?: number
}

/**
 * Cartoon avatar from owner reference photos.
 * Full-body idle; head tracks pointer. Prefer this while iterating; GLB mirrors the same mesh.
 */
export function StylizedHeroCharacter({ lookStrength = 0.35 }: StylizedHeroCharacterProps) {
  const group = useMemo(() => buildHeroCharacterGroup(), [])
  const pointer = useThree((state) => state.pointer)

  const parts = useMemo(() => {
    return {
      root: group,
      torso: group.getObjectByName('Torso') as THREE.Object3D | undefined,
      head: group.getObjectByName('Head') as THREE.Object3D | undefined,
      leftArm: group.getObjectByName('LeftArm') as THREE.Object3D | undefined,
      rightArm: group.getObjectByName('RightArm') as THREE.Object3D | undefined,
    }
  }, [group])

  const baseY = useRef(0)

  useFrame((state) => {
    const t = state.clock.elapsedTime
    const { root, torso, head, leftArm, rightArm } = parts
    root.position.y = baseY.current + Math.sin(t * 1.25) * 0.025
    root.rotation.y = Math.sin(t * 0.45) * 0.05
    if (torso) torso.scale.y = 1 + Math.sin(t * 1.25) * 0.01
    if (leftArm) {
      leftArm.rotation.x = Math.sin(t * 1.1) * 0.08 - 0.08
      leftArm.rotation.z = 0.12
    }
    if (rightArm) {
      rightArm.rotation.x = Math.sin(t * 1.1 + Math.PI) * 0.08 - 0.06
      rightArm.rotation.z = -0.12
    }
    if (head) {
      head.rotation.y = THREE.MathUtils.lerp(head.rotation.y, pointer.x * lookStrength, 0.08)
      head.rotation.x = THREE.MathUtils.lerp(
        head.rotation.x,
        -pointer.y * lookStrength * 0.4,
        0.08,
      )
    }
  })

  return <primitive object={group} />
}
